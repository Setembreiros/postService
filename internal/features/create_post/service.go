package create_post

import (
	"time"

	"github.com/rs/zerolog/log"
)

type CreatePostService struct {
}

type Post struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
}

func NewCreatePostService() *CreatePostService {
	return &CreatePostService{}
}

func (s *CreatePostService) CreatePost(post *Post) (string, string, error) {
	log.Info().Msgf("Post %s was created", post.Title)
	return "postId", "https://presigned/url", nil
}
