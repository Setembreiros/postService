package delete_post

import (
	database "postservice/internal/db"
	objectstorage "postservice/internal/objectStorage"

	"github.com/rs/zerolog/log"
)

type DeletePostRepository struct {
	dataRepository   *database.Database
	objectRepository *objectstorage.ObjectStorage
}

func NewDeletePostRepository(dataRepository *database.Database, objectRepository *objectstorage.ObjectStorage) *DeletePostRepository {
	return &DeletePostRepository{
		dataRepository:   dataRepository,
		objectRepository: objectRepository,
	}
}

func (r *DeletePostRepository) DeletePosts(postIds []string) error {
	posts, err := r.dataRepository.Client.GetPostsByIds(postIds)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error getting post metadatas for postIds %v", postIds)
		return err
	}

	var objectKeys []string
	for _, post := range posts {
		objectKey := post.User + "/" + post.Type + "/" + post.PostId + "." + post.FileType
		thumbnailObjectKey := post.User + "/" + post.Type + "/THUMBNAILS/" + post.PostId + "." + post.FileType
		objectKeys = append(objectKeys, objectKey)
		objectKeys = append(objectKeys, thumbnailObjectKey)
	}

	err = r.objectRepository.Client.DeleteObjects(objectKeys)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error deletting posts for postIds %v", postIds)
		return err
	}

	err = r.removePosts(postIds)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error deletting post metadatas for postIds %v", postIds)
		return err
	}

	return nil
}

func (r *DeletePostRepository) removePosts(postIds []string) error {
	postKeys := make([]any, len(postIds))
	for i, v := range postIds {
		postKeys[i] = &database.PostKey{
			PostId: v,
		}
	}

	err := r.dataRepository.Client.RemoveMultipleData("Posts", postKeys)

	return err
}
