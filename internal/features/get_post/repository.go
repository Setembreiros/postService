package get_post

import (
	database "postservice/internal/db"
	objectstorage "postservice/internal/objectStorage"

	"github.com/rs/zerolog/log"
)

type GetPostRepository struct {
	dataRepository   *database.Database
	objectRepository *objectstorage.ObjectStorage
}

func NewGetPostRepository(dataRepository *database.Database, objectRepository *objectstorage.ObjectStorage) *GetPostRepository {
	return &GetPostRepository{
		dataRepository:   dataRepository,
		objectRepository: objectRepository,
	}
}

func (r *GetPostRepository) GetPresignedUrlsForDownloading(username, lastPostId, lastPostCreatedAt string, limit int) ([]PostUrl, string, string, error) {
	posts, lastPostId, lastPostCreatedAt, err := r.getPostMetadatas(username, lastPostId, lastPostCreatedAt, limit)
	if err != nil {
		return []PostUrl{}, "", "", err
	}

	postUrls := r.getPostPresignedUrls(posts)

	return postUrls, lastPostId, lastPostCreatedAt, nil
}

func (r *GetPostRepository) getPostMetadatas(username, lastPostId, lastPostCreatedAt string, limit int) ([]*database.Post, string, string, error) {
	posts, lastPostId, lastPostCreatedAt, err := r.dataRepository.Client.GetPostsByIndexUser(username, lastPostId, lastPostCreatedAt, limit)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error getting post metadatas for username %s", username)
	}

	return posts, lastPostId, lastPostCreatedAt, err
}

func (r *GetPostRepository) getPostPresignedUrls(posts []*database.Post) []PostUrl {
	var postUrls []PostUrl

	for _, post := range posts {
		postUrl, err := r.getPostUrl(post)
		if err != nil {
			log.Error().Stack().Err(err).Msgf("Error getting presigned URLs for Post %s", post.PostId)
			continue
		}
		postUrls = append(postUrls, postUrl)
	}

	return postUrls
}

func (r *GetPostRepository) getPostUrl(post *database.Post) (PostUrl, error) {
	var postUrl PostUrl

	url, err := r.getPresignedUrl(post)
	if err != nil {
		return postUrl, err
	}
	thumbnailUrl := r.getPresignedThumbnailUrlIfExists(post)

	postUrl = PostUrl{
		PostId:                post.PostId,
		PresignedUrl:          url,
		PresignedThumbnailUrl: thumbnailUrl,
	}
	return postUrl, err
}

func (r *GetPostRepository) getPresignedUrl(post *database.Post) (string, error) {
	key := post.User + "/" + post.Type + "/" + post.PostId

	url, err := r.objectRepository.Client.GetPreSignedUrlForGettingObject(key)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error getting presigned URLs for Post %s", post.PostId)
		return "", err
	}

	return url, err
}

func (r *GetPostRepository) getPresignedThumbnailUrlIfExists(post *database.Post) string {
	if !post.HasThumbnail {
		return ""
	}

	thumbnailKey := post.User + "/" + post.Type + "/THUMBNAILS/" + post.PostId

	thumbnailUrl, err := r.objectRepository.Client.GetPreSignedUrlForGettingObject(thumbnailKey)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error getting presigned thumbnail URLs for Post %s", post.PostId)
		return ""
	}

	return thumbnailUrl
}
