package get_post_test

import (
	"bytes"
	"errors"
	"net/http/httptest"
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
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
	expectedPresignedUrls := []get_post.PostUrl{
		{
			PostId:       "post1",
			PresignedUrl: "url1",
		},
		{
			PostId:       "post2",
			PresignedUrl: "url2",
		},
		{
			PostId:       "post3",
			PresignedUrl: "url3",
		},
	}
	controllerRepository.EXPECT().GetPresignedUrlsForDownloading(username).Return(expectedPresignedUrls, nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {"urlPosts":[{"postId":"post1","url":"url1"},{"postId":"post2","url":"url2"},{"postId":"post3","url":"url3"}]}
	}`

	controller.GetUserPosts(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestInternalServerErrorOnGetUserPosts(t *testing.T) {
	setUpHandler(t)
	username := "username1"
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
	expectedError := errors.New("some error")
	controllerRepository.EXPECT().GetPresignedUrlsForDownloading(username).Return([]get_post.PostUrl{}, expectedError)
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError.Error() + `",
		"content":null
	}`

	controller.GetUserPosts(ginContext)

	assert.Equal(t, apiResponse.Code, 500)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func removeSpace(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\t", ""), "\n", "")
}
