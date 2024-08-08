package create_post_test

import (
	"bytes"
	"errors"
	"postservice/internal/features/create_post"
	mock_create_post "postservice/internal/features/create_post/mock"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var serviceLoggerOutput bytes.Buffer
var serviceRepository *mock_create_post.MockRepository
var createPostService *create_post.CreatePostService

func setUpService(t *testing.T) {
	ctrl := gomock.NewController(t)
	serviceRepository = mock_create_post.NewMockRepository(ctrl)
	log.Logger = log.Output(&serviceLoggerOutput)
	createPostService = create_post.NewCreatePostService(serviceRepository)
}

func TestCreatePostWithService(t *testing.T) {
	setUpService(t)
	newPost := &create_post.Post{
		User:        "username1",
		Title:       "Meu Post",
		Description: "Este é o meu novo post",
		CreatedAt:   time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
		LastUpdated: time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
	}
	expectedPostId := "username1-Meu_Post-1723153880"
	serviceRepository.EXPECT().AddNewPostMetaData(expectedPostId, newPost).Return(nil)

	postId, presignedUrl, err := createPostService.CreatePost(newPost)

	assert.Equal(t, postId, expectedPostId)
	assert.Equal(t, presignedUrl, "https://presigned/url")
	assert.Nil(t, err)
	assert.Contains(t, serviceLoggerOutput.String(), "Post Meu Post was created")
}

func TestErrorOnCreatePostWithService(t *testing.T) {
	setUpService(t)
	newPost := &create_post.Post{
		Title:       "Meu Post",
		Description: "Este é o meu novo post",
		CreatedAt:   time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
		LastUpdated: time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
	}
	expectedPostId := "username1-Meu_Post-1723153880"
	serviceRepository.EXPECT().AddNewPostMetaData(expectedPostId, newPost).Return(errors.New("some error"))

	postId, presignedUrl, err := createPostService.CreatePost(newPost)

	assert.Empty(t, postId)
	assert.Empty(t, presignedUrl)
	assert.NotNil(t, err)
	assert.Contains(t, serviceLoggerOutput.String(), "Error saving Post metadata")
}
