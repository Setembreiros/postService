package get_post_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"postservice/internal/features/get_post"
	mock_get_post "postservice/internal/features/get_post/mock"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
)

var controllerLoggerOutput bytes.Buffer
var controllerRepository *mock_get_post.MockRepository
var controller *get_post.GetPostController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func setUpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	controllerRepository = mock_get_post.NewMockRepository(ctrl)
	log.Logger = log.Output(&controllerLoggerOutput)
	controller = get_post.NewGetPostController(controllerRepository)
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)
}

func TestGetUserPost(t *testing.T) {
	setUpHandler(t)
	username := "username1"
	lastPostId := "post4"
	lastPostCreatedAt := "0001-01-03T00:00:00Z"
	limit := "4"
	ginContext.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	ginContext.Request.Method = "GET"
	ginContext.Request.Header.Set("Content-Type", "application/json")
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
	u := url.Values{}
	u.Add("lastPostId", lastPostId)
	u.Add("lastPostCreatedAt", lastPostCreatedAt)
	u.Add("limit", limit)
	ginContext.Request.URL.RawQuery = u.Encode()
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
	controllerRepository.EXPECT().GetPresignedUrlsForDownloading(username, lastPostId, lastPostCreatedAt, 4).Return(expectedPresignedUrls, "post7", "0001-01-06T00:00:00Z", nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {"urlPosts":[{"postId":"post1","url":"url1","thumbnailUrl":"thumbnailUrl1"},{"postId":"post2","url":"url2","thumbnailUrl":"thumbnailUrl2"},{"postId":"post3","url":"url3","thumbnailUrl":"thumbnailUrl3"}],"limit":4,"lastPostId":"post7","lastPostCreatedAt":"0001-01-06T00:00:00Z"}
	}`

	controller.GetUserPosts(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestGetUserPostWithDefaultPaginationParameters(t *testing.T) {
	setUpHandler(t)
	username := "username1"
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
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
	expectedDefaultLastPostId := ""
	expectedDefaultLastPostCreatedAt := ""
	expectedDefaultLimit := 6
	controllerRepository.EXPECT().GetPresignedUrlsForDownloading(username, expectedDefaultLastPostId, expectedDefaultLastPostCreatedAt, expectedDefaultLimit).Return(expectedPresignedUrls, "post7", "0001-01-06T00:00:00Z", nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {"urlPosts":[{"postId":"post1","url":"url1","thumbnailUrl":"thumbnailUrl1"},{"postId":"post2","url":"url2","thumbnailUrl":"thumbnailUrl2"},{"postId":"post3","url":"url3","thumbnailUrl":"thumbnailUrl3"}],"limit":6,"lastPostId":"post7","lastPostCreatedAt":"0001-01-06T00:00:00Z"}
	}`

	controller.GetUserPosts(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestInternalServerErrorOnGetUserPosts(t *testing.T) {
	setUpHandler(t)
	username := "username1"
	lastPostId := "post4"
	lastPostCreatedAt := "0001-01-03T00:00:00Z"
	limit := "4"
	ginContext.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	ginContext.Request.Method = "GET"
	ginContext.Request.Header.Set("Content-Type", "application/json")
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
	u := url.Values{}
	u.Add("lastPostId", lastPostId)
	u.Add("lastPostCreatedAt", lastPostCreatedAt)
	u.Add("limit", limit)
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedError := errors.New("some error")
	controllerRepository.EXPECT().GetPresignedUrlsForDownloading(username, lastPostId, lastPostCreatedAt, 4).Return([]get_post.PostUrl{}, "", "", expectedError)
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError.Error() + `",
		"content":null
	}`

	controller.GetUserPosts(ginContext)

	assert.Equal(t, apiResponse.Code, 500)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequestErrorOnGetUserPostsWhenLimitSmallerThanZero(t *testing.T) {
	setUpHandler(t)
	username := "username1"
	lastPostId := "post4"
	lastPostCreatedAt := "0001-01-03T00:00:00Z"
	wrongLimit := "0"
	ginContext.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	ginContext.Request.Method = "GET"
	ginContext.Request.Header.Set("Content-Type", "application/json")
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
	u := url.Values{}
	u.Add("lastPostId", lastPostId)
	u.Add("lastPostCreatedAt", lastPostCreatedAt)
	u.Add("limit", wrongLimit)
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedError := "Invalid pagination parameters, limit has to be greater than 0"
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError + `",
		"content":null
	}`

	controller.GetUserPosts(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequestErrorOnGetUserPostsWhenLastPostIdIsNotEmptyButLastPostCreatedAtIsEmpty(t *testing.T) {
	setUpHandler(t)
	username := "username1"
	lastPostId := "post4"
	lastPostCreatedAt := ""
	wrongLimit := "2"
	ginContext.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	ginContext.Request.Method = "GET"
	ginContext.Request.Header.Set("Content-Type", "application/json")
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
	u := url.Values{}
	u.Add("lastPostId", lastPostId)
	u.Add("lastPostCreatedAt", lastPostCreatedAt)
	u.Add("limit", wrongLimit)
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedError := "Invalid pagination parameters, lastPostId and lastPostCreatedAt both have to have value or both have to be empty"
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError + `",
		"content":null
	}`

	controller.GetUserPosts(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func removeSpace(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\t", ""), "\n", "")
}
