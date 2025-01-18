package create_post_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"postservice/internal/bus"
	mock_bus "postservice/internal/bus/mock"
	"postservice/internal/features/create_post"
	mock_create_post "postservice/internal/features/create_post/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var serviceLoggerOutput bytes.Buffer
var serviceRepository *mock_create_post.MockRepository
var serviceExternalBus *mock_bus.MockExternalBus
var serviceBus *bus.EventBus
var createPostService *create_post.CreatePostService

func setUpService(t *testing.T) {
	ctrl := gomock.NewController(t)
	serviceRepository = mock_create_post.NewMockRepository(ctrl)
	log.Logger = log.Output(&serviceLoggerOutput)
	serviceExternalBus = mock_bus.NewMockExternalBus(ctrl)
	serviceBus = bus.NewEventBus(serviceExternalBus)
	createPostService = create_post.NewCreatePostService(serviceRepository, serviceBus)
}

func TestCreatePostWithService(t *testing.T) {
	setUpService(t)
	newPost := &create_post.Post{
		User:         "username1",
		Type:         "Text",
		Title:        "Meu Post",
		Description:  "Este é o meu novo post",
		HasThumbnail: true,
	}
	serviceRepository.EXPECT().AddNewPostMetaData(newPost).Return(nil)
	serviceRepository.EXPECT().GetPresignedUrlsForUploading(newPost).Return([]string{"NoUploadId", "https://presigned/url", "https://presignedThumbanail/url"}, nil)

	result, err := createPostService.CreatePost(newPost)

	assert.Contains(t, result.PostId, "username1-Meu_Post-")
	assert.Contains(t, result.UploadId, "NoUploadId")
	assert.Equal(t, "https://presigned/url", result.PresignedUrls[0])
	assert.Equal(t, "https://presignedThumbanail/url", result.PresignedUrls[1])
	assert.Nil(t, err)
	assert.Contains(t, serviceLoggerOutput.String(), "Post Meu Post was created")
}

func TestErrorOnCreatePostWithService(t *testing.T) {
	setUpService(t)
	newPost := &create_post.Post{
		User:        "username1",
		Title:       "Meu Post",
		Type:        "Text",
		Description: "Este é o meu novo post",
	}
	serviceRepository.EXPECT().AddNewPostMetaData(newPost).Return(errors.New("some error"))
	serviceRepository.EXPECT().GetPresignedUrlsForUploading(newPost)

	result, err := createPostService.CreatePost(newPost)

	assert.Empty(t, result.PostId)
	assert.Empty(t, result.UploadId)
	assert.Empty(t, result.PresignedUrls)
	assert.NotNil(t, err)
	assert.Contains(t, serviceLoggerOutput.String(), "Error saving Post metadata")
}

func TestConfirmCreatedPostWithServiceWhenIsConfirmed(t *testing.T) {
	setUpService(t)
	postId := "postId"
	confirmedPost := &create_post.ConfirmedCreatedPost{
		IsConfirmed: true,
		PostId:      postId,
	}
	postMetadata := &create_post.Post{
		User:        "username1",
		Title:       "Meu Post",
		Type:        "Text",
		Description: "Este é o meu novo post",
	}
	expectedPostWasCreatedEvent := &create_post.PostWasCreatedEvent{
		PostId:   postId,
		Metadata: postMetadata,
	}
	expectedEvent, _ := createEvent("PostWasCreatedEvent", expectedPostWasCreatedEvent)
	serviceRepository.EXPECT().GetPostMetadata(postId).Return(postMetadata, nil)
	serviceExternalBus.EXPECT().Publish(expectedEvent)

	err := createPostService.ConfirmCreatedPost(confirmedPost)

	assert.Nil(t, err)
	assert.Contains(t, serviceLoggerOutput.String(), "Created Post postId was confirmed")
}

func TestErrorOnConfirmCreatedPostWithServiceWhenIsConfirmed(t *testing.T) {
	setUpService(t)
	postId := "postId"
	notConfirmedPost := &create_post.ConfirmedCreatedPost{
		IsConfirmed: true,
		PostId:      postId,
	}
	serviceRepository.EXPECT().GetPostMetadata(postId).Return(nil, errors.New("some error"))

	err := createPostService.ConfirmCreatedPost(notConfirmedPost)

	assert.NotNil(t, err)
	assert.Contains(t, serviceLoggerOutput.String(), "Error retrieving Post postId metadata")
}

func TestErrorOnConfirmCreatedPostWithServiceWhenSendingEvent(t *testing.T) {
	setUpService(t)
	postId := "postId"
	notConfirmedPost := &create_post.ConfirmedCreatedPost{
		IsConfirmed: true,
		PostId:      postId,
	}
	postMetadata := &create_post.Post{
		User:        "username1",
		Title:       "Meu Post",
		Type:        "Text",
		Description: "Este é o meu novo post",
	}
	expectedPostWasCreatedEvent := &create_post.PostWasCreatedEvent{
		PostId:   postId,
		Metadata: postMetadata,
	}
	expectedEvent, _ := createEvent("PostWasCreatedEvent", expectedPostWasCreatedEvent)
	serviceRepository.EXPECT().GetPostMetadata(postId).Return(postMetadata, nil)
	serviceExternalBus.EXPECT().Publish(expectedEvent).Return(errors.New("some error"))

	err := createPostService.ConfirmCreatedPost(notConfirmedPost)

	assert.NotNil(t, err)
	assert.Contains(t, serviceLoggerOutput.String(), "Publishing PostWasCreatedEvent failed")
}

func TestConfirmCreatedPostWithServiceWhenIsConfirmedAndIsMultipart(t *testing.T) {
	setUpService(t)
	postId := "postId"
	confirmedPost := &create_post.ConfirmedCreatedPost{
		IsConfirmed: true,
		PostId:      postId,
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
	postMetadata := &create_post.Post{
		User:        "username1",
		Title:       "Meu Post",
		Type:        "Text",
		Description: "Este é o meu novo post",
	}
	expectedMultipartPost := &create_post.MultipartPost{
		Post:           postMetadata,
		UploadId:       confirmedPost.UploadId,
		CompletedParts: confirmedPost.CompletedParts,
	}
	expectedPostWasCreatedEvent := &create_post.PostWasCreatedEvent{
		PostId:   postId,
		Metadata: postMetadata,
	}
	expectedEvent, _ := createEvent("PostWasCreatedEvent", expectedPostWasCreatedEvent)
	serviceRepository.EXPECT().GetPostMetadata(postId).Return(postMetadata, nil)
	serviceRepository.EXPECT().CompleteMultipartUpload(expectedMultipartPost).Return(nil)
	serviceExternalBus.EXPECT().Publish(expectedEvent)

	err := createPostService.ConfirmCreatedPost(confirmedPost)

	assert.Nil(t, err)
	assert.Contains(t, serviceLoggerOutput.String(), "Created Post postId was confirmed")
}

func TestErrorOnConfirmCreatedPostWithServiceWhenCompleteMultipartUpload(t *testing.T) {
	setUpService(t)
	postId := "postId"
	confirmedPost := &create_post.ConfirmedCreatedPost{
		IsConfirmed: true,
		PostId:      postId,
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
	postMetadata := &create_post.Post{
		User:        "username1",
		Title:       "Meu Post",
		Type:        "Text",
		Description: "Este é o meu novo post",
	}
	expectedMultipartPost := &create_post.MultipartPost{
		Post:           postMetadata,
		UploadId:       confirmedPost.UploadId,
		CompletedParts: confirmedPost.CompletedParts,
	}
	serviceRepository.EXPECT().GetPostMetadata(postId).Return(postMetadata, nil)
	serviceRepository.EXPECT().CompleteMultipartUpload(expectedMultipartPost).Return(errors.New("some error"))

	err := createPostService.ConfirmCreatedPost(confirmedPost)

	assert.NotNil(t, err)
	assert.Contains(t, serviceLoggerOutput.String(), "Error completing multipart Post")
}

func TestConfirmCreatedPostWithServiceWhenIsNotConfirmed(t *testing.T) {
	setUpService(t)
	notConfirmedPost := &create_post.ConfirmedCreatedPost{
		IsConfirmed: false,
		PostId:      "postId",
	}
	serviceRepository.EXPECT().RemoveUnconfirmedPost(notConfirmedPost.PostId).Return(nil)

	err := createPostService.ConfirmCreatedPost(notConfirmedPost)

	assert.Nil(t, err)
	assert.Contains(t, serviceLoggerOutput.String(), "Created Post postId failed")
}

func TestErrorOnConfirmCreatedPostWithServiceWhenIsNotConfirmed(t *testing.T) {
	setUpService(t)
	notConfirmedPost := &create_post.ConfirmedCreatedPost{
		IsConfirmed: false,
		PostId:      "postId",
	}
	serviceRepository.EXPECT().RemoveUnconfirmedPost(notConfirmedPost.PostId).Return(errors.New("some error"))

	err := createPostService.ConfirmCreatedPost(notConfirmedPost)

	assert.NotNil(t, err)
	assert.Contains(t, serviceLoggerOutput.String(), "Error removing Post metadata")
}

func createEvent(eventName string, eventData any) (*bus.Event, error) {
	dataEvent, err := serialize(eventData)
	if err != nil {
		return nil, err
	}

	return &bus.Event{
		Type: eventName,
		Data: dataEvent,
	}, nil
}

func serialize(data any) ([]byte, error) {
	return json.Marshal(data)
}
