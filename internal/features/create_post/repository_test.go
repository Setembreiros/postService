package create_post_test

import (
	database "postservice/internal/db"
	mock_database "postservice/internal/db/mock"
	"postservice/internal/features/create_post"
	objectstorage "postservice/internal/objectStorage"
	mock_objectstorage "postservice/internal/objectStorage/mock"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

var dbClient *mock_database.MockDatabaseClient
var osClient *mock_objectstorage.MockObjectStorageClient
var createPostRepository *create_post.CreatePostRepository

func setUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbClient = mock_database.NewMockDatabaseClient(ctrl)
	osClient = mock_objectstorage.NewMockObjectStorageClient(ctrl)
	createPostRepository = create_post.NewCreatePostRepository(database.NewDatabase(dbClient), objectstorage.NewObjectStorage(osClient))
}

func TestAddNewPostMetaDataInRepository(t *testing.T) {
	setUp(t)
	newPost := &create_post.Post{
		PostId:       "username1-Meu_Post-1723153880",
		User:         "username1",
		Type:         "Text",
		Title:        "Meu Post",
		Description:  "Este é o meu novo post",
		HasThumbnail: true,
		CreatedAt:    time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC().String(),
		LastUpdated:  time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC().String(),
	}
	data := &create_post.PostMetadata{
		PostId:       newPost.PostId,
		User:         newPost.User,
		Type:         newPost.Type,
		Title:        newPost.Title,
		Description:  newPost.Description,
		HasThumbnail: newPost.HasThumbnail,
		CreatedAt:    newPost.CreatedAt,
		LastUpdated:  newPost.LastUpdated,
	}
	dbClient.EXPECT().InsertData("Posts", data)

	createPostRepository.AddNewPostMetaData(newPost)
}

func TestGetPresignedUrlsForUploading_HasThumbnailIsTrue(t *testing.T) {
	setUp(t)
	newPost := &create_post.Post{
		PostId:       "username1-Meu_Post-1723153880",
		User:         "username1",
		Type:         "Text",
		Title:        "Meu Post",
		Description:  "Este é o meu novo post",
		Size:         50,
		HasThumbnail: true,
		CreatedAt:    time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC().String(),
		LastUpdated:  time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC().String(),
	}
	expectedKey := "username1/Text/username1-Meu_Post-1723153880"
	expectedThumbnailKey := "username1/Text/THUMBNAILS/username1-Meu_Post-1723153880"
	osClient.EXPECT().GetPreSignedUrlsForPuttingObject(expectedKey, newPost.Size)
	osClient.EXPECT().GetPreSignedUrlsForPuttingObject(expectedThumbnailKey, 0).Return("NoUploadId", []string{"fakeurl"}, nil)

	createPostRepository.GetPresignedUrlsForUploading(newPost)
}

func TestGetPresignedUrlsForUploading_HasThumbnailIsFalse(t *testing.T) {
	setUp(t)
	newPost := &create_post.Post{
		PostId:       "username1-Meu_Post-1723153880",
		User:         "username1",
		Type:         "Text",
		Title:        "Meu Post",
		Description:  "Este é o meu novo post",
		Size:         50,
		HasThumbnail: false,
		CreatedAt:    time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC().String(),
		LastUpdated:  time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC().String(),
	}
	expectedKey := "username1/Text/username1-Meu_Post-1723153880"
	osClient.EXPECT().GetPreSignedUrlsForPuttingObject(expectedKey, newPost.Size)

	createPostRepository.GetPresignedUrlsForUploading(newPost)
}

func TestGetPostMetadata(t *testing.T) {
	setUp(t)
	var post create_post.Post
	postId := "username1-Meu_Post-1723153880"
	expectedKey := &create_post.PostKey{
		PostId: postId,
	}
	dbClient.EXPECT().GetData("Posts", expectedKey, &post)

	createPostRepository.GetPostMetadata(postId)
}

func TestRemoveUnconfirmedPostMetaDataInRepository(t *testing.T) {
	setUp(t)
	postId := "username1-Meu_Post-1723153880"
	expectedKey := &create_post.PostKey{
		PostId: postId,
	}
	dbClient.EXPECT().RemoveData("Posts", expectedKey)

	createPostRepository.RemoveUnconfirmedPost(postId)
}
