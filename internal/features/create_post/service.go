package create_post

import (
	"time"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=mock/service.go

type Repository interface {
	AddNewPostMetaData(data *Post) error
	GetPresignedUrlForUploadingText(data *Post) (string, error)
	RemoveUnconfirmedPost(postId string) error
}

type CreatePostService struct {
	repository Repository
}

type Post struct {
	User        string    `json:"username"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
}

type ConfirmedCreatedPost struct {
	IsConfirmed bool   `json:"is_confirmed"`
	PostId      string `json:"post_id"`
}

func NewCreatePostService(repository Repository) *CreatePostService {
	return &CreatePostService{
		repository: repository,
	}
}

func (s *CreatePostService) CreatePost(post *Post) (string, error) {
	chError := make(chan error, 2)
	chResult := make(chan string, 1)

	go s.savePostMetaData(post, chError)
	go s.generetePreSignedUrl(post, chResult, chError)

	numberOfTasks := 2
	for i := 0; i < numberOfTasks; i++ {
		err := <-chError
		if err != nil {
			return "", err
		}
	}

	result := <-chResult
	log.Info().Msgf("Post %s was created", post.Title)
	return result, nil
}

func (s *CreatePostService) ConfirmCreatedPost(post *ConfirmedCreatedPost) error {
	if post.IsConfirmed {
		log.Info().Msgf("Created Post %s was confirmed", post.PostId)
	} else {
		err := s.repository.RemoveUnconfirmedPost(post.PostId)
		if err != nil {
			log.Error().Stack().Err(err).Msg("Error removing Post metadata")
			return err
		}
		log.Info().Msgf("Created Post %s failed", post.PostId)
	}
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

func (s *CreatePostService) generetePreSignedUrl(post *Post, chResult chan string, chError chan<- error) {
	presignedUrl, err := s.repository.GetPresignedUrlForUploadingText(post)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error generating Pre-Signed URL")
		chError <- err
	}

	chError <- nil
	chResult <- presignedUrl
}
