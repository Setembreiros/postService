package delete_post_test

import (
	"bytes"
	"errors"
	"fmt"
	database "postservice/internal/db"
	mock_database "postservice/internal/db/mock"
	"postservice/internal/features/delete_post"
	objectStorage "postservice/internal/objectStorage"
	mock_objectStorage "postservice/internal/objectStorage/mock"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var repositoryLoggerOutput bytes.Buffer
var dataClient *mock_database.MockDatabaseClient
var objectClient *mock_objectStorage.MockObjectStorageClient
var deletePostRepository *delete_post.DeletePostRepository

func setUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	log.Logger = log.Output(&repositoryLoggerOutput)
	dataClient = mock_database.NewMockDatabaseClient(ctrl)
	objectClient = mock_objectStorage.NewMockObjectStorageClient(ctrl)
	deletePostRepository = delete_post.NewDeletePostRepository(database.NewDatabase(dataClient), objectStorage.NewObjectStorage(objectClient))
}

func TestDeletePostsWithRepository(t *testing.T) {
	setUp(t)
	postIds := []string{"1", "2", "3"}
	data := []*database.Post{
		{
			PostId:    "usernam1-meuPost-170948521",
			User:      "usernam1",
			Title:     "meuPost",
			FileType:  "png",
			CreatedAt: time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
		},
		{
			PostId:    "usernam1-meuPost2-184639321",
			User:      "usernam1",
			Title:     "meuPost2",
			FileType:  "pdf",
			CreatedAt: time.Date(2024, 7, 24, 20, 51, 20, 33, time.UTC).UTC(),
		},
	}
	expectedKeys := []string{
		data[0].User + "/" + data[0].Type + "/" + data[0].PostId + "." + data[0].FileType,
		data[0].User + "/" + data[0].Type + "/THUMBNAILS/" + data[0].PostId + "." + data[0].FileType,
		data[1].User + "/" + data[1].Type + "/" + data[1].PostId + "." + data[1].FileType,
		data[1].User + "/" + data[1].Type + "/THUMBNAILS/" + data[1].PostId + "." + data[1].FileType,
	}
	expectedPostKeys := []any{
		&database.PostKey{
			PostId: "1",
		},
		&database.PostKey{
			PostId: "2",
		},
		&database.PostKey{
			PostId: "3",
		},
	}
	dataClient.EXPECT().GetPostsByIds(postIds).Return(data, nil)
	objectClient.EXPECT().DeleteObjects(expectedKeys)
	dataClient.EXPECT().RemoveMultipleData("Posts", expectedPostKeys).Return(nil)

	deletePostRepository.DeletePosts(postIds)
}

func TestDeletePostsWithRepository_GettingPostMetadataError(t *testing.T) {
	setUp(t)
	postIds := []string{"1", "2", "3"}
	dataClient.EXPECT().GetPostsByIds(postIds).Return(nil, errors.New("some error"))

	deletePostRepository.DeletePosts(postIds)

	assert.Contains(t, repositoryLoggerOutput.String(), fmt.Sprintf("Error getting post metadatas for postIds %v", postIds))
}

func TestDeletePostsWithRepository_DeletingObjectsError(t *testing.T) {
	setUp(t)
	postIds := []string{"1", "2", "3"}
	data := []*database.Post{
		{
			PostId:    "usernam1-meuPost-170948521",
			User:      "usernam1",
			Title:     "meuPost",
			FileType:  "png",
			CreatedAt: time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
		},
		{
			PostId:    "usernam1-meuPost2-184639321",
			User:      "usernam1",
			Title:     "meuPost2",
			FileType:  "pdf",
			CreatedAt: time.Date(2024, 7, 24, 20, 51, 20, 33, time.UTC).UTC(),
		},
	}
	expectedKeys := []string{
		data[0].User + "/" + data[0].Type + "/" + data[0].PostId + "." + data[0].FileType,
		data[0].User + "/" + data[0].Type + "/THUMBNAILS/" + data[0].PostId + "." + data[0].FileType,
		data[1].User + "/" + data[1].Type + "/" + data[1].PostId + "." + data[1].FileType,
		data[1].User + "/" + data[1].Type + "/THUMBNAILS/" + data[1].PostId + "." + data[1].FileType,
	}
	dataClient.EXPECT().GetPostsByIds(postIds).Return(data, nil)
	objectClient.EXPECT().DeleteObjects(expectedKeys).Return(errors.New("some error"))

	deletePostRepository.DeletePosts(postIds)

	assert.Contains(t, repositoryLoggerOutput.String(), fmt.Sprintf("Error deletting posts for postIds %v", postIds))
}

func TestDeletePostsWithRepository_RemovingPostMetadataError(t *testing.T) {
	setUp(t)
	postIds := []string{"1", "2", "3"}
	data := []*database.Post{
		{
			PostId:    "usernam1-meuPost-170948521",
			User:      "usernam1",
			Title:     "meuPost",
			FileType:  "png",
			CreatedAt: time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
		},
		{
			PostId:    "usernam1-meuPost2-184639321",
			User:      "usernam1",
			Title:     "meuPost2",
			FileType:  "pdf",
			CreatedAt: time.Date(2024, 7, 24, 20, 51, 20, 33, time.UTC).UTC(),
		},
	}
	expectedKeys := []string{
		data[0].User + "/" + data[0].Type + "/" + data[0].PostId + "." + data[0].FileType,
		data[0].User + "/" + data[0].Type + "/THUMBNAILS/" + data[0].PostId + "." + data[0].FileType,
		data[1].User + "/" + data[1].Type + "/" + data[1].PostId + "." + data[1].FileType,
		data[1].User + "/" + data[1].Type + "/THUMBNAILS/" + data[1].PostId + "." + data[1].FileType,
	}
	expectedPostKeys := []any{
		&database.PostKey{
			PostId: "1",
		},
		&database.PostKey{
			PostId: "2",
		},
		&database.PostKey{
			PostId: "3",
		},
	}
	dataClient.EXPECT().GetPostsByIds(postIds).Return(data, nil)
	objectClient.EXPECT().DeleteObjects(expectedKeys)
	dataClient.EXPECT().RemoveMultipleData("Posts", expectedPostKeys).Return(errors.New("some error"))

	deletePostRepository.DeletePosts(postIds)

	assert.Contains(t, repositoryLoggerOutput.String(), fmt.Sprintf("Error deletting post metadatas for postIds %v", postIds))
}
