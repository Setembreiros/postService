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
	GetPresignedUrlsForUploading(data *Post) (PresignedUrl, error)
	GetPostMetadata(postId string) (*Post, error)
	CompleteMultipartUpload(multipartPost *MultipartPost) error
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

type CreatePostResult struct {
	PostId       string       `json:"postId"`
	PresignedUrl PresignedUrl `json:"presignedUrl"`
}

type ConfirmedCreatedPost struct {
	IsConfirmed    bool            `json:"isConfirmed"`
	PostId         string          `json:"postId"`
	IsMultipart    bool            `json:"isMultipart"`
	UploadId       string          `json:"uploadId"`
	CompletedParts []CompletedPart `json:"completedParts"`
}

type MultipartPost struct {
	Post           *Post           `json:"post"`
	UploadId       string          `json:"uploadId"`
	CompletedParts []CompletedPart `json:"completedParts"`
}

type CompletedPart struct {
	PartNumber int    `json:"partNumber"`
	ETag       string `json:"eTag"`
}

type PostWasCreatedEvent struct {
	PostId   string `json:"post_id"`
	Metadata *Post  `json:"metadata"`
}

type PresignedUrl struct {
	UploadId              string   `json:"uploadId"`
	ContentPresignedUrls  []string `json:"contentPresignedUrls"`
	ThumbanilPresignedUrl string   `json:"thumbanilPresignedUrl"`
}

func NewCreatePostService(repository Repository, bus *bus.EventBus) *CreatePostService {
	return &CreatePostService{
		repository: repository,
		bus:        bus,
	}
}

var timeLayout string = "2006-01-02T15:04:05.000000Z"

func (s *CreatePostService) CreatePost(post *Post) (CreatePostResult, error) {
	chError := make(chan error, 2)
	chResult := make(chan PresignedUrl, 1)

	post.CreatedAt = time.Now().UTC().Format(timeLayout)
	post.LastUpdated = post.CreatedAt
	postId, err := generatePostId(post)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error generating Post Id")
		return CreatePostResult{}, err
	}
	post.PostId = postId

	go s.savePostMetaData(post, chError)
	go s.generetePreSignedUrl(post, chResult, chError)

	numberOfTasks := 2
	for i := 0; i < numberOfTasks; i++ {
		err := <-chError
		if err != nil {
			return CreatePostResult{}, err
		}
	}

	result := <-chResult
	log.Info().Msgf("Post %s was created", post.Title)
	return CreatePostResult{
		PostId:       postId,
		PresignedUrl: result,
	}, nil
}

func (s *CreatePostService) ConfirmCreatedPost(confirmPostData *ConfirmedCreatedPost) error {
	if !confirmPostData.IsConfirmed {
		err := s.rollBackUnconfirmedPost(confirmPostData.PostId)
		if err != nil {
			return err
		}

		return nil
	}

	post, err := s.repository.GetPostMetadata(confirmPostData.PostId)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error retrieving Post %s metadata", confirmPostData.PostId)
		return err
	}

	if confirmPostData.IsMultipart {
		multipartPost := &MultipartPost{
			Post:           post,
			UploadId:       confirmPostData.UploadId,
			CompletedParts: confirmPostData.CompletedParts,
		}
		err := s.repository.CompleteMultipartUpload(multipartPost)
		log.Error().Stack().Err(err).Msgf("Error completing multipart Post %s", confirmPostData.PostId)
		if err != nil {
			return err
		}
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

func (s *CreatePostService) generetePreSignedUrl(post *Post, chResult chan PresignedUrl, chError chan<- error) {
	presignedUrl, err := s.repository.GetPresignedUrlsForUploading(post)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error generating Pre-Signed URL")
		chError <- err
	}

	chError <- nil
	chResult <- presignedUrl
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
