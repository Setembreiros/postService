package get_post

import (
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=mock/service.go

type Repository interface {
	GetPresignedUrlsForDownloading(username string) ([]string, error)
}

type GetPostService struct {
	repository Repository
}

func NewGetPostService(repository Repository) *GetPostService {
	return &GetPostService{
		repository: repository,
	}
}

func (s *GetPostService) GetUserPosts(username string) ([]string, error) {
	presignedUrls, err := s.repository.GetPresignedUrlsForDownloading(username)
	if err != nil {
		return presignedUrls, err
	}

	log.Info().Msgf("%s's Pre-Signed Url Posts were generated", username)
	return presignedUrls, nil
}
