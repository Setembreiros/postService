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

func (r *GetPostRepository) GetPresignedUrlsForDownloading(username string) ([]PostUrl, error) {
	posts, err := r.getPostMetadatas(username)
	if err != nil {
		return []PostUrl{}, err
	}

	postUrls := r.getPostPresignedUrls(posts)

	return postUrls, nil
}

func (r *GetPostRepository) getPostMetadatas(username string) ([]*database.Post, error) {
	posts, err := r.dataRepository.Client.GetPostsByIndexUser(username)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error getting post metadatas for username %s", username)
	}

	return posts, err
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
