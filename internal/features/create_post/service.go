package create_post

import (
	"postservice/internal/bus"
	"time"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=mock/service.go

type Repository interface {
	AddNewPostMetaData(data *Post) error
	GetPresignedUrlsForUploading(data *Post) ([]string, error)
	GetPostMetadata(postId string) (*Post, error)
	RemoveUnconfirmedPost(postId string) error
}

type CreatePostService struct {
	repository Repository
	bus        *bus.EventBus
}

type Post struct {
	User         string    `json:"username"`
	Type         string    `json:"type"`
	FileType     string    `json:"fileType"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	HasThumbnail bool      `json:"hasThumbnail"`
	CreatedAt    time.Time `json:"createdAt"`
	LastUpdated  time.Time `json:"lastUpdated"`
}

type ConfirmedCreatedPost struct {
	IsConfirmed bool   `json:"is_confirmed"`
	PostId      string `json:"post_id"`
}

type PostWasCreatedEvent struct {
	PostId   string `json:"post_id"`
	Metadata *Post  `json:"metadata"`
}

func NewCreatePostService(repository Repository, bus *bus.EventBus) *CreatePostService {
	return &CreatePostService{
		repository: repository,
		bus:        bus,
	}
}

func (s *CreatePostService) CreatePost(post *Post) (string, []string, error) {
	chError := make(chan error, 2)
	chResult := make(chan []string, 1)

	go s.savePostMetaData(post, chError)
	go s.generetePreSignedUrl(post, chResult, chError)

	numberOfTasks := 2
	for i := 0; i < numberOfTasks; i++ {
		err := <-chError
		if err != nil {
			return "", []string{}, err
		}
	}

	result := <-chResult
	log.Info().Msgf("Post %s was created", post.Title)
	return generatePostId(post), result, nil
}

func (s *CreatePostService) ConfirmCreatedPost(confirmPostData *ConfirmedCreatedPost) error {
	if !confirmPostData.IsConfirmed {
		err := s.rollBackUnconfirmedPost(confirmPostData.PostId)
		if err != nil {
			return err
		}

		return nil
	}

	post, err := s.repository.GetPostMetadata(confirmPostData.PostId)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error retrieving Post %s metadata", confirmPostData.PostId)
		return err
	}

	err = s.publishPostWasCreatedEvent(confirmPostData.PostId, post)
	if err != nil {
		return err
	}

	log.Info().Msgf("Created Post %s was confirmed", confirmPostData.PostId)

	return nil
}

func (s *CreatePostService) savePostMetaData(post *Post, chError chan<- error) {
	err := s.repository.AddNewPostMetaData(post)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error saving Post metadata")
		chError <- err
	}

	chError <- nil
}

func (s *CreatePostService) generetePreSignedUrl(post *Post, chResult chan []string, chError chan<- error) {
	presignedUrls, err := s.repository.GetPresignedUrlsForUploading(post)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error generating Pre-Signed URL")
		chError <- err
	}

	chError <- nil
	chResult <- presignedUrls
}

func (s *CreatePostService) publishPostWasCreatedEvent(postId string, metadata *Post) error {
	event := &PostWasCreatedEvent{
		PostId:   postId,
		Metadata: metadata,
	}
	err := s.bus.Publish("PostWasCreatedEvent", event)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Publishing PostWasCreatedEvent failed")
		return err
	}

	return nil
}

func (s *CreatePostService) rollBackUnconfirmedPost(postId string) error {
	err := s.repository.RemoveUnconfirmedPost(postId)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error removing Post metadata")
		return err
	}
	log.Info().Msgf("Created Post %s failed", postId)

	return nil
}
