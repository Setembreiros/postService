package delete_post_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"postservice/internal/bus"
	mock_bus "postservice/internal/bus/mock"
	database "postservice/internal/db"
	"postservice/internal/features/delete_post"
	mock_delete_post "postservice/internal/features/delete_post/mock"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
)

var controllerLoggerOutput bytes.Buffer
var controllerRepository *mock_delete_post.MockRepository
var controllerExternalBus *mock_bus.MockExternalBus
var controllerBus *bus.EventBus
var controller *delete_post.DeletePostController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func setUpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	controllerRepository = mock_delete_post.NewMockRepository(ctrl)
	log.Logger = log.Output(&controllerLoggerOutput)
	controllerExternalBus = mock_bus.NewMockExternalBus(ctrl)
	controllerBus = bus.NewEventBus(controllerExternalBus)
	controller = delete_post.NewDeletePostController(controllerRepository, controllerBus)
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)
}

func TestDeletePosts(t *testing.T) {
	setUpHandler(t)
	username := "username1"
	req, _ := http.NewRequest("DELETE", "/posts/username1?postId=1&postId=2&postId=3", nil)
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
	ginContext.Request = req
	controllerRepository.EXPECT().DeletePosts([]string{"1", "2", "3"}).Return(nil)
	expectedPostsWereDeletedEvent := &delete_post.PostsWereDeletedEvent{
		Username: username,
		PostIds:  []string{"1", "2", "3"},
	}
	expectedEvent := createEvent("PostsWereDeletedEvent", expectedPostsWereDeletedEvent)
	controllerExternalBus.EXPECT().Publish(expectedEvent).Return(nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": null
	}`

	controller.DeletePosts(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestDeletePosts_MissingPostID(t *testing.T) {
	setUpHandler(t)
	req, _ := http.NewRequest("DELETE", "/posts/username1", nil)
	ginContext.Request = req
	expectedBodyResponse := `{
		"error": true,
		"message": "Missing postId parameters",
		"content": null
	}`

	controller.DeletePosts(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestDeletePosts_NotFound(t *testing.T) {
	setUpHandler(t)
	req, _ := http.NewRequest("DELETE", "/posts/username1?postId=1&postId=2&postId=3", nil)
	ginContext.Request = req
	controllerRepository.EXPECT().DeletePosts([]string{"1", "2", "3"}).Return(database.NewNotFoundError("", "2"))
	expectedBodyResponse := `{
		"error": true,
		"message": "` + fmt.Sprintf("Some posts were not found for post ids %v", []string{"1", "2", "3"}) + `",
		"content": null
	}`

	controller.DeletePosts(ginContext)

	assert.Equal(t, apiResponse.Code, 404)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestDeletePosts_InternalServerError(t *testing.T) {
	setUpHandler(t)
	req, _ := http.NewRequest("DELETE", "/posts/username1?postId=1&postId=2&postId=3", nil)
	ginContext.Request = req
	controllerRepository.EXPECT().DeletePosts([]string{"1", "2", "3"}).Return(errors.New("Some error"))
	expectedBodyResponse := `{
		"error": true,
		"message": "Some error",
		"content": null
	}`

	controller.DeletePosts(ginContext)

	assert.Equal(t, apiResponse.Code, 500)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func removeSpace(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\t", ""), "\n", "")
}
