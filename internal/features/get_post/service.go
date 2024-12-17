package get_post

import (
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=mock/service.go

type Repository interface {
	GetPresignedUrlsForDownloading(username, lastPostId, lastPostCreatedAt string, limit int) ([]PostUrl, string, string, error)
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

func (s *GetPostService) GetUserPosts(username, lastPostId, lastPostCreatedAt string, limit int) ([]PostUrl, string, string, error) {
	postUrls, lastPostId, lastPostCreatedAt, err := s.repository.GetPresignedUrlsForDownloading(username, lastPostId, lastPostCreatedAt, limit)
	if err != nil {
		return postUrls, lastPostId, lastPostCreatedAt, err
	}

	log.Info().Msgf("%s's Pre-Signed Url Posts were generated", username)
	return postUrls, lastPostId, lastPostCreatedAt, nil
}
