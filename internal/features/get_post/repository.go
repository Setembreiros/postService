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
	posts, err := r.dataRepository.Client.GetPostsByIndexUser(username)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error getting post metadatas for username %s", username)
		return []PostUrl{}, err
	}

	var postUrls []PostUrl
	for _, post := range posts {
		key := post.User + "/" + post.Type + "/" + post.PostId + "." + post.FileType
		url, err := r.objectRepository.Client.GetPreSignedUrlForGettingObject(key)
		posturl := PostUrl{
			PostId:       post.PostId,
			PresignedUrl: url,
		}
		if err != nil {
			log.Error().Stack().Err(err).Msgf("Error getting presigned URLs for Post %s", post.PostId)
			continue
		}
		postUrls = append(postUrls, posturl)
	}

	return postUrls, nil
}
