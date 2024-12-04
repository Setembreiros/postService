package get_post

import (
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=mock/service.go

type Repository interface {
	GetPresignedUrlsForDownloading(username, lastPostId string, limit int) ([]PostUrl, string, error)
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

func (s *GetPostService) GetUserPosts(username, lastPostId string, limit int) ([]PostUrl, string, error) {
	postUrls, nextPostId, err := s.repository.GetPresignedUrlsForDownloading(username, lastPostId, limit)
	if err != nil {
		return postUrls, nextPostId, err
	}

	log.Info().Msgf("%s's Pre-Signed Url Posts were generated", username)
	return postUrls, nextPostId, nil
}
