package create_post_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"postservice/internal/features/create_post"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/rs/zerolog/log"
)

var controllerLoggerOutput bytes.Buffer
var controller *create_post.CreatePostController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func setUpHandler(t *testing.T) {
	log.Logger = log.Output(&controllerLoggerOutput)
	controller = create_post.NewCreatePostController()
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)
}

func TestCreatePost(t *testing.T) {
	setUpHandler(t)
	newPost := &create_post.Post{
		Title:       "Meu Post",
		Description: "Este é o meu novo post",
		CreatedAt:   time.Now(),
		LastUpdated: time.Now(),
	}
	data, _ := serializeData(newPost)
	ginContext.Request = httptest.NewRequest(http.MethodPost, "/post", bytes.NewBuffer(data))
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"post_id":"postId",
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
