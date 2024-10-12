package get_post_test

import (
	"bytes"
	"errors"
	database "postservice/internal/db"
	mock_database "postservice/internal/db/mock"
	"postservice/internal/features/get_post"
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
var getPostRepository *get_post.GetPostRepository

func setUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	log.Logger = log.Output(&repositoryLoggerOutput)
	dataClient = mock_database.NewMockDatabaseClient(ctrl)
	objectClient = mock_objectStorage.NewMockObjectStorageClient(ctrl)
	getPostRepository = get_post.NewGetPostRepository(database.NewDatabase(dataClient), objectStorage.NewObjectStorage(objectClient))
}

func TestGetPresignedUrlsForDownloadingInRepository(t *testing.T) {
	setUp(t)
	username := "username1"
	data := []*database.Post{
		{
			PostId:       "usernam1-meuPost-170948521",
			User:         username,
			Title:        "meuPost",
			HasThumbnail: true,
			CreatedAt:    time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
		},
		{
			PostId:       "usernam1-meuPost2-184639321",
			User:         username,
			Title:        "meuPost2",
			HasThumbnail: true,
			CreatedAt:    time.Date(2024, 7, 24, 20, 51, 20, 33, time.UTC).UTC(),
		},
		{
			PostId:       "usernam1-meuPost3-184639321",
			User:         username,
			Title:        "meuPost3",
			HasThumbnail: false,
			CreatedAt:    time.Date(2024, 7, 24, 20, 51, 20, 33, time.UTC).UTC(),
		},
	}
	expectedKey1 := data[0].User + "/" + data[0].Type + "/" + data[0].PostId
	expectedThumbnailKey1 := data[0].User + "/" + data[0].Type + "/THUMBNAILS/" + data[0].PostId
	expectedKey2 := data[1].User + "/" + data[1].Type + "/" + data[1].PostId
	expectedThumbnailKey2 := data[1].User + "/" + data[1].Type + "/THUMBNAILS/" + data[1].PostId
	expectedKey3 := data[2].User + "/" + data[2].Type + "/" + data[2].PostId
	dataClient.EXPECT().GetPostsByIndexUser(username).Return(data, nil)
	objectClient.EXPECT().GetPreSignedUrlForGettingObject(expectedKey1)
	objectClient.EXPECT().GetPreSignedUrlForGettingObject(expectedThumbnailKey1)
	objectClient.EXPECT().GetPreSignedUrlForGettingObject(expectedKey2)
	objectClient.EXPECT().GetPreSignedUrlForGettingObject(expectedThumbnailKey2)
	objectClient.EXPECT().GetPreSignedUrlForGettingObject(expectedKey3)

	getPostRepository.GetPresignedUrlsForDownloading(username)
}

func TestErrorOnGetPresignedUrlsForDownloadingInRepositoryWhenGettingPostMetadataByIndexuser(t *testing.T) {
	setUp(t)
	username := "username1"
	dataClient.EXPECT().GetPostsByIndexUser(username).Return(nil, errors.New("some error"))

	getPostRepository.GetPresignedUrlsForDownloading(username)

	assert.Contains(t, repositoryLoggerOutput.String(), "Error getting post metadatas for username "+username)
}

func TestErrorOnGetPresignedUrlsForDownloadingInRepositoryWhenGettingUrls(t *testing.T) {
	setUp(t)
	username := "username1"
	data := []*database.Post{
		{
			PostId:       "usernam1-meuPost-170948521",
			User:         username,
			Title:        "meuPost",
			HasThumbnail: true,
			CreatedAt:    time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
		},
		{
			PostId:       "usernam1-meuPost2-184639321",
			User:         username,
			Title:        "meuPost2",
			HasThumbnail: true,
			CreatedAt:    time.Date(2024, 7, 24, 20, 51, 20, 33, time.UTC).UTC(),
		},
	}
	expectedKey1 := data[0].User + "/" + data[0].Type + "/" + data[0].PostId
	expectedKey2 := data[1].User + "/" + data[1].Type + "/" + data[1].PostId
	expectedThumbnailKey2 := data[1].User + "/" + data[1].Type + "/THUMBNAILS/" + data[1].PostId
	expectedResult := []get_post.PostUrl{
		{
			PostId:                "usernam1-meuPost2-184639321",
			PresignedUrl:          "url2",
			PresignedThumbnailUrl: "thumbnailUrl2",
		},
	}
	dataClient.EXPECT().GetPostsByIndexUser(username).Return(data, nil)
	objectClient.EXPECT().GetPreSignedUrlForGettingObject(expectedKey1).Return("", errors.New("some error"))
	objectClient.EXPECT().GetPreSignedUrlForGettingObject(expectedKey2).Return(expectedResult[0].PresignedUrl, nil)
	objectClient.EXPECT().GetPreSignedUrlForGettingObject(expectedThumbnailKey2).Return(expectedResult[0].PresignedThumbnailUrl, nil)

	result, err := getPostRepository.GetPresignedUrlsForDownloading(username)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, expectedResult, result)
	assert.Contains(t, repositoryLoggerOutput.String(), "Error getting presigned URLs for Post "+data[0].PostId)
}
