package get_post_test

import (
	"bytes"
	"errors"
	"postservice/internal/features/get_post"
	mock_get_post "postservice/internal/features/get_post/mock"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var serviceLoggerOutput bytes.Buffer
var serviceRepository *mock_get_post.MockRepository
var getPostService *get_post.GetPostService

func setUpService(t *testing.T) {
	ctrl := gomock.NewController(t)
	serviceRepository = mock_get_post.NewMockRepository(ctrl)
	log.Logger = log.Output(&serviceLoggerOutput)
	getPostService = get_post.NewGetPostService(serviceRepository)
}

func TestGetUserPostsWithService(t *testing.T) {
	setUpService(t)
	username := "username1"
	expectedPostMetadatas := []*get_post.Post{
		{
			User:        username,
			Type:        "TEXT",
			FileType:    "pdf",
			Title:       "Meu Post",
			Description: "Este é o meu novo post",
			CreatedAt:   time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
			LastUpdated: time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
		},
		{
			User:        username,
			Type:        "IMAGE",
			FileType:    "png",
			Title:       "Meu Post",
			Description: "Este é o meu novo post",
			CreatedAt:   time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
			LastUpdated: time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
		},
		{
			User:        username,
			Type:        "VIDEO",
			FileType:    "mp4",
			Title:       "Meu Post",
			Description: "Este é o meu novo post",
			CreatedAt:   time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
			LastUpdated: time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
		},
	}
	expectedPresignedUrls := []string{"url1", "url2", "url3"}
	serviceRepository.EXPECT().GetUserPostMetadatas(username).Return(expectedPostMetadatas, nil)
	serviceRepository.EXPECT().GetPresignedUrlsForDownloading(expectedPostMetadatas).Return(expectedPresignedUrls, nil)

	getPostService.GetUserPosts(username)

	assert.Contains(t, serviceLoggerOutput.String(), username+"'s Pre-Signed Url Posts were generated")
}

func TestErrorOnGetUserPostsWithServiceWhenGettingMetadata(t *testing.T) {
	setUpService(t)
	username := "username1"
	expectedPostMetadatas := []*get_post.Post{}
	serviceRepository.EXPECT().GetUserPostMetadatas(username).Return(expectedPostMetadatas, errors.New("some error"))

	getPostService.GetUserPosts(username)

	assert.Contains(t, serviceLoggerOutput.String(), "Error getting post metadatas for username "+username)
}

func TestErrorOnGetUserPostsWithServiceWhenGettingUrls(t *testing.T) {
	setUpService(t)
	username := "username1"
	expectedPostMetadatas := []*get_post.Post{
		{
			User:        username,
			Type:        "TEXT",
			FileType:    "pdf",
			Title:       "Meu Post",
			Description: "Este é o meu novo post",
			CreatedAt:   time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
			LastUpdated: time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
		},
		{
			User:        username,
			Type:        "IMAGE",
			FileType:    "png",
			Title:       "Meu Post",
			Description: "Este é o meu novo post",
			CreatedAt:   time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
			LastUpdated: time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
		},
		{
			User:        username,
			Type:        "VIDEO",
			FileType:    "mp4",
			Title:       "Meu Post",
			Description: "Este é o meu novo post",
			CreatedAt:   time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
			LastUpdated: time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
		},
	}
	expectedPresignedUrls := []string{}
	serviceRepository.EXPECT().GetUserPostMetadatas(username).Return(expectedPostMetadatas, nil)
	serviceRepository.EXPECT().GetPresignedUrlsForDownloading(expectedPostMetadatas).Return(expectedPresignedUrls, errors.New("some error"))

	getPostService.GetUserPosts(username)

	assert.Contains(t, serviceLoggerOutput.String(), "Error getting presigned URLs for username "+username)
}
