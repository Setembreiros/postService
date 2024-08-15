package get_post

import (
	"time"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=mock/service.go

type Repository interface {
	GetUserPostMetadatas(username string) ([]*Post, error)
	GetPresignedUrlsForDownloading(data []*Post) ([]string, error)
}

type GetPostService struct {
	repository Repository
}

type Post struct {
	User        string    `json:"username"`
	Type        string    `json:"type"`
	FileType    string    `json:"file_type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
}

type ConfirmedCreatedPost struct {
	IsConfirmed bool   `json:"is_confirmed"`
	PostId      string `json:"post_id"`
}

type PostWasCreatedEvent struct {
	PostId   string `json:"post_id"`
	Metadata *Post  `json:"metadata"`
}

func NewGetPostService(repository Repository) *GetPostService {
	return &GetPostService{
		repository: repository,
	}
}

func (s *GetPostService) GetUserPosts(username string) ([]string, error) {
	posts, err := s.repository.GetUserPostMetadatas(username)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error getting post metadatas for username %s", username)
		return []string{}, err
	}

	presignedUrls, err := s.repository.GetPresignedUrlsForDownloading(posts)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error getting presigned URLs for username %s", username)
		return presignedUrls, err
	}

	log.Info().Msgf("%s's Pre-Signed Url Posts were generated", username)
	return presignedUrls, nil
}
