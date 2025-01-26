package create_post_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"postservice/internal/bus"
	"postservice/internal/features/create_post"
	mock_create_post "postservice/internal/features/create_post/mock"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
)

var controllerLoggerOutput bytes.Buffer
var controllerService *mock_create_post.MockService
var controllerBus *bus.EventBus
var controller *create_post.CreatePostController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func setUpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	controllerService = mock_create_post.NewMockService(ctrl)
	controllerBus = &bus.EventBus{}
	log.Logger = log.Output(&controllerLoggerOutput)
	controller = create_post.NewCreatePostController(controllerService, controllerBus)
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)
}

func TestCreatePost_HasThumbnailIsTrue(t *testing.T) {
	setUpHandler(t)
	newPost := &create_post.Post{
		User:         "username1",
		Type:         "Text",
		Title:        "Meu Post",
		Description:  "Este é o meu novo post",
		Size:         120,
		HasThumbnail: true,
	}
	data, _ := serializeData(newPost)
	ginContext.Request = httptest.NewRequest(http.MethodPost, "/post", bytes.NewBuffer(data))
	expectedPostId := "username1-Meu_Post-1723153880"
	expectedPresignedUrl1 := "https://presigned/url1"
	expectedPresignedUrl2 := "https://presigned/url2"
	expectedPresignedUrlThumbanil := "https://presigned/url/thumbnail"
	controllerService.EXPECT().CreatePost(newPost).Return(create_post.CreatePostResult{expectedPostId, create_post.PresignedUrl{"NoUploadId", []string{expectedPresignedUrl1, expectedPresignedUrl2}, expectedPresignedUrlThumbanil}}, nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"postId": "` + expectedPostId + `",
			"uploadId": "NoUploadId",
			"presignedUrls":["` + expectedPresignedUrl1 + `","` + expectedPresignedUrl2 + `"],
			"presignedThumbnailUrl":"` + expectedPresignedUrlThumbanil + `"
		}
	}`

	controller.CreatePost(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestCreatePost_HasThumbnailIsFalse(t *testing.T) {
	setUpHandler(t)
	newPost := &create_post.Post{
		User:         "username1",
		Type:         "Text",
		Title:        "Meu Post",
		Description:  "Este é o meu novo post",
		Size:         120,
		HasThumbnail: false,
	}
	data, _ := serializeData(newPost)
	ginContext.Request = httptest.NewRequest(http.MethodPost, "/post", bytes.NewBuffer(data))
	expectedPostId := "username1-Meu_Post-1723153880"
	expectedPresignedUrl1 := "https://presigned/url1"
	expectedPresignedUrl2 := "https://presigned/url2"
	controllerService.EXPECT().CreatePost(newPost).Return(create_post.CreatePostResult{expectedPostId, create_post.PresignedUrl{"NoUploadId", []string{expectedPresignedUrl1, expectedPresignedUrl2}, ""}}, nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"postId": "` + expectedPostId + `",
			"uploadId": "NoUploadId",
			"presignedUrls":["` + expectedPresignedUrl1 + `","` + expectedPresignedUrl2 + `"],
			"presignedThumbnailUrl":""
		}
	}`

	controller.CreatePost(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestInternalServerErrorOnCreatePost(t *testing.T) {
	setUpHandler(t)
	newPost := &create_post.Post{
		User:        "username1",
		Type:        "Text",
		Title:       "Meu Post",
		Description: "Este é o meu novo post",
	}
	data, _ := serializeData(newPost)
	ginContext.Request = httptest.NewRequest(http.MethodPost, "/post", bytes.NewBuffer(data))
	expectedError := errors.New("some error")
	controllerService.EXPECT().CreatePost(newPost).Return(create_post.CreatePostResult{}, expectedError)
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError.Error() + `",
		"content":null
	}`

	controller.CreatePost(ginContext)

	assert.Equal(t, apiResponse.Code, 500)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestConfirmCreatedPostWhenIsConfirmed(t *testing.T) {
	setUpHandler(t)
	confirmedPost := &create_post.ConfirmedCreatedPost{
		IsConfirmed: true,
		PostId:      "postId",
	}
	data, _ := serializeData(confirmedPost)
	ginContext.Request = httptest.NewRequest(http.MethodPut, "/confirm-created-post", bytes.NewBuffer(data))
	controllerService.EXPECT().ConfirmCreatedPost(confirmedPost)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": null
	}`

	controller.ConfirmCreatedPost(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestConfirmCreatedPostWhenMultipartIsConfirmed(t *testing.T) {
	setUpHandler(t)
	confirmedPost := &create_post.ConfirmedCreatedPost{
		IsConfirmed: true,
		PostId:      "postId",
		IsMultipart: true,
		UploadId:    "upload-id",
		CompletedParts: []create_post.CompletedPart{
			{
				PartNumber: 1,
				ETag:       "etag1",
			},
			{
				PartNumber: 2,
				ETag:       "etag2",
			},
		},
	}
	data, _ := serializeData(confirmedPost)
	ginContext.Request = httptest.NewRequest(http.MethodPut, "/confirm-created-post", bytes.NewBuffer(data))
	controllerService.EXPECT().ConfirmCreatedPost(confirmedPost)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": null
	}`

	controller.ConfirmCreatedPost(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestConfirmCreatedPostWhenIsNotConfirmed(t *testing.T) {
	setUpHandler(t)
	notConfirmedPost := &create_post.ConfirmedCreatedPost{
		IsConfirmed: false,
		PostId:      "postId",
	}
	data, _ := serializeData(notConfirmedPost)
	ginContext.Request = httptest.NewRequest(http.MethodPut, "/confirm-created-post", bytes.NewBuffer(data))
	controllerService.EXPECT().ConfirmCreatedPost(notConfirmedPost)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": null
	}`

	controller.ConfirmCreatedPost(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func serializeData(data any) ([]byte, error) {
	return json.Marshal(data)
}

func removeSpace(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\t", ""), "\n", "")
}
