package delete_post

import (
	"postservice/internal/bus"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=mock/service.go

type Repository interface {
	DeletePosts(postIds []string) error
}

type DeletePostService struct {
	repository Repository
	bus        *bus.EventBus
}

type PostsWereDeletedEvent struct {
	Username string   `json:"username"`
	PostIds  []string `json:"postIds"`
}

func NewDeletePostService(repository Repository, bus *bus.EventBus) *DeletePostService {
	return &DeletePostService{
		repository: repository,
		bus:        bus,
	}
}

func (s *DeletePostService) DeletePosts(username string, postIds []string) error {
	err := s.repository.DeletePosts(postIds)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error deleting posts for postIds %v", postIds)
		return err
	}

	s.publishPostsWereDeletedEvent(username, postIds)

	log.Info().Msgf("%v were deleted", postIds)
	return nil
}

func (s *DeletePostService) publishPostsWereDeletedEvent(username string, postIds []string) error {
	event := &PostsWereDeletedEvent{
		Username: username,
		PostIds:  postIds,
	}
	err := s.bus.Publish("PostsWereDeletedEvent", event)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Publishing PostsWereDeletedEvent failed")
		return err
	}

	return nil
}
