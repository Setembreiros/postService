package create_post

import (
	"postservice/internal/bus"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=mock/service.go

type Repository interface {
	AddNewPostMetaData(data *Post) error
	GetPresignedUrlsForUploading(data *Post) ([]string, error)
	GetPostMetadata(postId string) (*Post, error)
	CompleteMultipartUpload(multipartPost MultipartPost) error
	RemoveUnconfirmedPost(postId string) error
}

type CreatePostService struct {
	repository Repository
	bus        *bus.EventBus
}

type Post struct {
	PostId       string `json:"postId"`
	User         string `json:"username"`
	Type         string `json:"type"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Size         int    `json:"size"`
	HasThumbnail bool   `json:"hasThumbnail"`
	CreatedAt    string `json:"createdAt"`
	LastUpdated  string `json:"lastUpdated"`
}

type ConfirmedCreatedPost struct {
	IsConfirmed   bool          `json:"isConfirmed"`
	PostId        string        `json:"postId"`
	IsMultipart   bool          `json:"isMultipart"`
	MultipartPost MultipartPost `json:"multipartPost"`
}

type MultipartPost struct {
	Key           string          `json:"key"`
	UploadID      string          `json:"uploadID"`
	CompletedPart []CompletedPart `json:"completedPart"`
}

type CompletedPart struct {
	PartNumber int    `json:"partNumber"`
	ETag       string `json:"eTag"`
}

type PostWasCreatedEvent struct {
	PostId   string `json:"post_id"`
	Metadata *Post  `json:"metadata"`
}

func NewCreatePostService(repository Repository, bus *bus.EventBus) *CreatePostService {
	return &CreatePostService{
		repository: repository,
		bus:        bus,
	}
}

var timeLayout string = "2006-01-02T15:04:05.000000000Z"

func (s *CreatePostService) CreatePost(post *Post) (string, []string, error) {
	chError := make(chan error, 2)
	chResult := make(chan []string, 1)

	post.CreatedAt = time.Now().UTC().Format(timeLayout)
	post.LastUpdated = post.CreatedAt
	postId, err := generatePostId(post)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error generating Post Id")
		return "", []string{}, err
	}
	post.PostId = postId

	go s.savePostMetaData(post, chError)
	go s.generetePreSignedUrl(post, chResult, chError)

	numberOfTasks := 2
	for i := 0; i < numberOfTasks; i++ {
		err := <-chError
		if err != nil {
			return "", []string{}, err
		}
	}

	result := <-chResult
	log.Info().Msgf("Post %s was created", post.Title)
	return post.PostId, result, nil
}

func (s *CreatePostService) ConfirmCreatedPost(confirmPostData *ConfirmedCreatedPost) error {
	if !confirmPostData.IsConfirmed {
		err := s.rollBackUnconfirmedPost(confirmPostData.PostId)
		if err != nil {
			return err
		}

		return nil
	}

	if confirmPostData.IsMultipart {
		err := s.repository.CompleteMultipartUpload(confirmPostData.MultipartPost)
		log.Error().Stack().Err(err).Msgf("Error completing multipart Post %s", confirmPostData.PostId)
		if err != nil {
			return err
		}
	}

	post, err := s.repository.GetPostMetadata(confirmPostData.PostId)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error retrieving Post %s metadata", confirmPostData.PostId)
		return err
	}

	err = s.publishPostWasCreatedEvent(confirmPostData.PostId, post)
	if err != nil {
		return err
	}

	log.Info().Msgf("Created Post %s was confirmed", confirmPostData.PostId)

	return nil
}

func (s *CreatePostService) savePostMetaData(post *Post, chError chan<- error) {
	err := s.repository.AddNewPostMetaData(post)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error saving Post metadata")
		chError <- err
	}

	chError <- nil
}

func (s *CreatePostService) generetePreSignedUrl(post *Post, chResult chan []string, chError chan<- error) {
	presignedUrls, err := s.repository.GetPresignedUrlsForUploading(post)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error generating Pre-Signed URL")
		chError <- err
	}

	chError <- nil
	chResult <- presignedUrls
}

func (s *CreatePostService) publishPostWasCreatedEvent(postId string, metadata *Post) error {
	event := &PostWasCreatedEvent{
		PostId:   postId,
		Metadata: metadata,
	}

	err := s.bus.Publish("PostWasCreatedEvent", event)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Publishing PostWasCreatedEvent failed")
		return err
	}

	return nil
}

func (s *CreatePostService) rollBackUnconfirmedPost(postId string) error {
	err := s.repository.RemoveUnconfirmedPost(postId)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error removing Post metadata")
		return err
	}
	log.Info().Msgf("Created Post %s failed", postId)

	return nil
}

func generatePostId(post *Post) (string, error) {
	parsedCreatedAt, err := time.Parse(timeLayout, post.CreatedAt)
	if err != nil {
		return "", err
	}
	postId := post.User + "-" + post.Title + "-" + strconv.FormatInt(parsedCreatedAt.Unix(), 10)
	return strings.ReplaceAll(strings.ReplaceAll(postId, " ", "_"), "\t", "_"), nil
}
