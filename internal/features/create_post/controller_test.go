package create_post_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"postservice/internal/features/create_post"
	"strings"
	"testing"

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
	data, _ := serializeData("")
	ginContext.Request = httptest.NewRequest(http.MethodPost, "/post", bytes.NewBuffer(data))
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": "Everything Ok"
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
