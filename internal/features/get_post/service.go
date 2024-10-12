package get_post

import (
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=mock/service.go

type Repository interface {
	GetPresignedUrlsForDownloading(username string) ([]PostUrl, error)
}

type GetPostService struct {
	repository Repository
}

type PostUrl struct {
	PostId                string `json:"postId"`
	PresignedUrl          string `json:"url"`
	PresignedThumbnailUrl string `json:"thumbnailUrl"`
}

func NewGetPostService(repository Repository) *GetPostService {
	return &GetPostService{
		repository: repository,
	}
}

func (s *GetPostService) GetUserPosts(username string) ([]PostUrl, error) {
	postUrls, err := s.repository.GetPresignedUrlsForDownloading(username)
	if err != nil {
		return postUrls, err
	}

	log.Info().Msgf("%s's Pre-Signed Url Posts were generated", username)
	return postUrls, nil
}
