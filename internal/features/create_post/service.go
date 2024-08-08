package create_post

import (
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=mock/service.go

type Repository interface {
	AddNewPostMetaData(id string, data *Post) error
}

type CreatePostService struct {
	repository Repository
}

type Post struct {
	User        string    `json:"username"`
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

func (s *CreatePostService) CreatePost(post *Post) (string, string, error) {
	postId := generatePostId(post)
	err := s.savePostMetaData(postId, post)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error saving Post metadata")
		return "", "", err
	}

	log.Info().Msgf("Post %s was created", post.Title)

	return postId, "https://presigned/url", nil
}

func generatePostId(post *Post) string {
	postId := post.User + "-" + post.Title + "-" + strconv.FormatInt(post.CreatedAt.Unix(), 10)
	return strings.ReplaceAll(strings.ReplaceAll(postId, " ", "_"), "\t", "_")
}

func (s *CreatePostService) savePostMetaData(id string, post *Post) error {
	err := s.repository.AddNewPostMetaData(id, post)
	return err
}

func (s *CreatePostService) GeneretePresSignedUrl(post *Post) (string, error) {
	log.Info().Msgf("Post %s was created", post.Title)
	return "https://presigned/url", nil
}
