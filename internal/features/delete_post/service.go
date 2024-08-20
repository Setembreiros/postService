package delete_post

import (
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=mock/service.go

type Repository interface {
	DeletePosts(postIds []string) error
}

type DeletePostService struct {
	repository Repository
}

func NewDeletePostService(repository Repository) *DeletePostService {
	return &DeletePostService{
		repository: repository,
	}
}

func (s *DeletePostService) DeletePosts(postIds []string) error {
	err := s.repository.DeletePosts(postIds)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error deleting posts for postIds %v", postIds)
		return err
	}

	log.Info().Msgf("%v were deleted", postIds)
	return nil
}
