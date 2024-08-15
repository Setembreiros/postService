package get_post_test

import (
	"bytes"
	"errors"
	"net/http/httptest"
	database "postservice/internal/db"
	"postservice/internal/features/get_post"
	mock_get_post "postservice/internal/features/get_post/mock"
	"strings"
	"testing"
	"time"

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
	controllerRepository.EXPECT().GetUserPostMetadatas(username).Return(expectedPostMetadatas, nil)
	controllerRepository.EXPECT().GetPresignedUrlsForDownloading(expectedPostMetadatas).Return(expectedPresignedUrls, nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": ["url1", "url2", "url3"]
	}`

	controller.GetUserPosts(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestUserNotFoundOnGetUserPosts(t *testing.T) {
	setUpHandler(t)
	noExistingUsername := "noExistingUsername"
	ginContext.Params = []gin.Param{{Key: "username", Value: noExistingUsername}}
	expectedPostMetadatas := []*get_post.Post{}
	expectedNotFoundError := &database.NotFoundError{}
	controllerRepository.EXPECT().GetUserPostMetadatas(noExistingUsername).Return(expectedPostMetadatas, expectedNotFoundError)
	expectedBodyResponse := `{
		"error": true,
		"message": "User not found for username ` + noExistingUsername + `",
		"content":null
	}`

	controller.GetUserPosts(ginContext)

	assert.Equal(t, apiResponse.Code, 404)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestInternalServerErrorOnGetUserPosts(t *testing.T) {
	setUpHandler(t)
	username := "username1"
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
	expectedPostMetadatas := []*get_post.Post{}
	expectedError := errors.New("some error")
	controllerRepository.EXPECT().GetUserPostMetadatas(username).Return(expectedPostMetadatas, expectedError)
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
