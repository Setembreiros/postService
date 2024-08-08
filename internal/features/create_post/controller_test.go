package create_post_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"postservice/internal/features/create_post"
	mock_create_post "postservice/internal/features/create_post/mock"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
)

var controllerLoggerOutput bytes.Buffer
var controllerRepository *mock_create_post.MockRepository
var controller *create_post.CreatePostController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func setUpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	controllerRepository = mock_create_post.NewMockRepository(ctrl)
	log.Logger = log.Output(&controllerLoggerOutput)
	controller = create_post.NewCreatePostController(controllerRepository)
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)
}

func TestCreatePost(t *testing.T) {
	setUpHandler(t)
	newPost := &create_post.Post{
		User:        "username1",
		Title:       "Meu Post",
		Description: "Este é o meu novo post",
		CreatedAt:   time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
		LastUpdated: time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
	}
	data, _ := serializeData(newPost)
	ginContext.Request = httptest.NewRequest(http.MethodPost, "/post", bytes.NewBuffer(data))
	expectedPostId := "username1-Meu_Post-1723153880"
	controllerRepository.EXPECT().AddNewPostMetaData(expectedPostId, newPost)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"post_id":"` + expectedPostId + `",
			"presigned_url":"https://presigned/url"
		}
	}`

	controller.CreatePost(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func serializeData(data any) ([]byte, error) {
	return json.Marshal(data)
}

func removeSpace(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\t", ""), "\n", "")
}
