package get_post_test

import (
	"bytes"
	"errors"
	"postservice/internal/features/get_post"
	mock_get_post "postservice/internal/features/get_post/mock"
	"testing"

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
	serviceLoggerOutput.Truncate(0)
	log.Logger = log.Output(&serviceLoggerOutput)
	getPostService = get_post.NewGetPostService(serviceRepository)
}

func TestGetUserPostsWithService(t *testing.T) {
	setUpService(t)
	username := "username1"
	lastPostId := "post"
	limit := 3
	expectedPresignedUrls := []get_post.PostUrl{
		{
			PostId:                "post1",
			PresignedUrl:          "url1",
			PresignedThumbnailUrl: "thumbnailUrl1",
		},
		{
			PostId:                "post2",
			PresignedUrl:          "url2",
			PresignedThumbnailUrl: "thumbnailUrl2",
		},
		{
			PostId:                "post3",
			PresignedUrl:          "url3",
			PresignedThumbnailUrl: "thumbnailUrl3",
		},
	}
	serviceRepository.EXPECT().GetPresignedUrlsForDownloading(username, lastPostId, limit).Return(expectedPresignedUrls, "post4", nil)

	getPostService.GetUserPosts(username, lastPostId, limit)

	assert.Contains(t, serviceLoggerOutput.String(), username+"'s Pre-Signed Url Posts were generated")
}

func TestErrorOnGetUserPostsWithServiceWhenGettingUrls(t *testing.T) {
	setUpService(t)
	username := "username1"
	lastPostId := "post"
	limit := 2
	serviceRepository.EXPECT().GetPresignedUrlsForDownloading(username, lastPostId, limit).Return(nil, "", errors.New("some error"))

	getPostService.GetUserPosts(username, lastPostId, limit)

	assert.NotContains(t, serviceLoggerOutput.String(), username+"'s Pre-Signed Url Posts were generated")
}
