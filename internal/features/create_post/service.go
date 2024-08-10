package create_post

import (
	"time"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=mock/service.go

type Repository interface {
	AddNewPostMetaData(data *Post) error
	GetPresignedUrlForUploadingText(data *Post) (string, error)
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

func NewCreatePostService(repository Repository) *CreatePostService {
	return &CreatePostService{
		repository: repository,
	}
}

func (s *CreatePostService) CreatePost(post *Post) (string, error) {
	err := s.savePostMetaData(post)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error saving Post metadata")
		return "", err
	}
	presignedUrl, err := s.generetePreSignedUrl(post)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error generating Pre-Signed URL")
		return "", err
	}

	log.Info().Msgf("Post %s was created", post.Title)

	return presignedUrl, nil
}

func (s *CreatePostService) savePostMetaData(post *Post) error {
	err := s.repository.AddNewPostMetaData(post)
	return err
}

func (s *CreatePostService) generetePreSignedUrl(post *Post) (string, error) {
	presignedUrl, err := s.repository.GetPresignedUrlForUploadingText(post)
	return presignedUrl, err
}
