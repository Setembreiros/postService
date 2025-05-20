package delete_post_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"postservice/internal/bus"
	mock_bus "postservice/internal/bus/mock"
	"postservice/internal/features/delete_post"
	mock_delete_post "postservice/internal/features/delete_post/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var serviceLoggerOutput bytes.Buffer
var serviceRepository *mock_delete_post.MockRepository
var serviceExternalBus *mock_bus.MockExternalBus
var serviceBus *bus.EventBus
var deletePostService *delete_post.DeletePostService

func setUpService(t *testing.T) {
	ctrl := gomock.NewController(t)
	serviceRepository = mock_delete_post.NewMockRepository(ctrl)
	log.Logger = log.Output(&serviceLoggerOutput)
	serviceExternalBus = mock_bus.NewMockExternalBus(ctrl)
	serviceBus = bus.NewEventBus(serviceExternalBus)
	deletePostService = delete_post.NewDeletePostService(serviceRepository, serviceBus)
}

func TestDeletePostsWithService(t *testing.T) {
	setUpService(t)
	username := "username1"
	postIds := []string{"1", "2", "3"}
	serviceRepository.EXPECT().DeletePosts(postIds).Return(nil)
	expectedPostsWereDeletedEvent := &delete_post.PostsWereDeletedEvent{
		Username: username,
		PostIds:  postIds,
	}
	expectedEvent := createEvent("PostsWereDeletedEvent", expectedPostsWereDeletedEvent)
	serviceExternalBus.EXPECT().Publish(expectedEvent)

	deletePostService.DeletePosts(username, postIds)

	assert.Contains(t, serviceLoggerOutput.String(), "[1 2 3] were deleted")
}

func TestDeletePostsWithService_Error(t *testing.T) {
	setUpService(t)
	username := "username1"
	postIds := []string{"1", "2", "3", "4"}
	serviceRepository.EXPECT().DeletePosts(postIds).Return(errors.New("Some error"))

	deletePostService.DeletePosts(username, postIds)

	assert.Contains(t, serviceLoggerOutput.String(), fmt.Sprintf("Error deleting posts for postIds %v", postIds))
}

func TestDeletePostsWithService_ErrorPublishingEvent(t *testing.T) {
	setUpService(t)
	username := "username1"
	postIds := []string{"1", "2", "3", "4"}
	serviceRepository.EXPECT().DeletePosts(postIds).Return(nil)
	expectedPostsWereDeletedEvent := &delete_post.PostsWereDeletedEvent{
		Username: username,
		PostIds:  postIds,
	}
	expectedEvent := createEvent("PostsWereDeletedEvent", expectedPostsWereDeletedEvent)
	serviceExternalBus.EXPECT().Publish(expectedEvent).Return(errors.New("some error"))

	deletePostService.DeletePosts(username, postIds)

	assert.Contains(t, serviceLoggerOutput.String(), "Publishing PostsWereDeletedEvent failed")
}

func createEvent(eventName string, eventData any) *bus.Event {
	dataEvent, err := serialize(eventData)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &bus.Event{
		Type: eventName,
		Data: dataEvent,
	}
}

func serialize(data any) ([]byte, error) {
	return json.Marshal(data)
}
