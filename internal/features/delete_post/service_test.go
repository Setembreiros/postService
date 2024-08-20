package delete_post_test

import (
	"bytes"
	"errors"
	"fmt"
	"postservice/internal/features/delete_post"
	mock_delete_post "postservice/internal/features/delete_post/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var serviceLoggerOutput bytes.Buffer
var serviceRepository *mock_delete_post.MockRepository
var deletePostService *delete_post.DeletePostService

func setUpService(t *testing.T) {
	ctrl := gomock.NewController(t)
	serviceRepository = mock_delete_post.NewMockRepository(ctrl)
	log.Logger = log.Output(&serviceLoggerOutput)
	deletePostService = delete_post.NewDeletePostService(serviceRepository)
}

func TestDeletePostsWithService(t *testing.T) {
	setUpService(t)
	postIds := []string{"1", "2", "3"}
	serviceRepository.EXPECT().DeletePosts(postIds).Return(nil)

	deletePostService.DeletePosts(postIds)

	assert.Contains(t, serviceLoggerOutput.String(), "[1 2 3] were deleted")
}

func TestDeletePostsWithService_Error(t *testing.T) {
	setUpService(t)
	postIds := []string{"1", "2", "3", "4"}
	serviceRepository.EXPECT().DeletePosts(postIds).Return(errors.New("Some error"))

	deletePostService.DeletePosts(postIds)

	assert.Contains(t, serviceLoggerOutput.String(), fmt.Sprintf("Error deleting posts for postIds %v", postIds))
}
