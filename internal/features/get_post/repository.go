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

func (r *GetPostRepository) GetPresignedUrlsForDownloading(username string) ([]string, error) {
	posts, err := r.dataRepository.Client.GetPostsByIndexUser(username)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error getting post metadatas for username %s", username)
		return []string{}, err
	}

	var presignedUrls []string
	for _, post := range posts {
		key := post.User + "/" + post.Type + "/" + post.PostId + "." + post.FileType
		url, err := r.objectRepository.Client.GetPreSignedUrlForPuttingObject(key)
		if err != nil {
			log.Error().Stack().Err(err).Msgf("Error getting presigned URLs for Post %s", post.PostId)
			continue
		}
		presignedUrls = append(presignedUrls, url)
	}

	return presignedUrls, nil
}
