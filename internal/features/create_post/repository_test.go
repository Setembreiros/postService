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
		User:        "username1",
		Type:        "Text",
		Title:       "Meu Post",
		Description: "Este é o meu novo post",
		CreatedAt:   time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
		LastUpdated: time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
	}
	expectedPostId := "username1-Meu_Post-1723153880"
	data := &create_post.PostMetadata{
		PostId:   expectedPostId,
		Metadata: newPost,
	}
	dbClient.EXPECT().InsertData("Posts", data)

	createPostRepository.AddNewPostMetaData(newPost)
}

func TestGetPresignedUrlForUploadingText(t *testing.T) {
	setUp(t)
	newPost := &create_post.Post{
		User:        "username1",
		Type:        "Text",
		Title:       "Meu Post",
		Description: "Este é o meu novo post",
		CreatedAt:   time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
		LastUpdated: time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
	}
	expectedKey := "username1/Text/username1-Meu_Post-1723153880"
	osClient.EXPECT().GetPreSignedUrlForPuttingObject(expectedKey)

	createPostRepository.GetPresignedUrlForUploadingText(newPost)
}
