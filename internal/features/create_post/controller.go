package create_post

import (
	"postservice/internal/api"
	"postservice/internal/bus"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=controller.go -destination=mock/controller.go

type CreatePostController struct {
	service Service
}

type CreatePostResponse struct {
	PostId                string   `json:"postId"`
	UploadId              string   `json:"uploadId"`
	PresignedUrls         []string `json:"presignedUrls"`
	PresignedThumbnailUrl string   `json:"presignedThumbnailUrl"`
}

type Service interface {
	CreatePost(post *Post) (CreatePostResult, error)
	ConfirmCreatedPost(confirmPostData *ConfirmedCreatedPost) error
}

func NewCreatePostController(service Service, bus *bus.EventBus) *CreatePostController {
	return &CreatePostController{
		service: service,
	}
}

func (controller *CreatePostController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/post", controller.CreatePost)
	routerGroup.PUT("/confirm-created-post", controller.ConfirmCreatedPost)
}

func (controller *CreatePostController) CreatePost(c *gin.Context) {
	log.Info().Msg("Handling Request POST CreatePost")
	var post Post

	if err := c.BindJSON(&post); err != nil {
		log.Error().Stack().Err(err).Msg("Invalid Data")
		api.SendBadRequest(c, "Invalid Json Request")
		return
	}

	postResult, err := controller.service.CreatePost(&post)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	postResponse := &CreatePostResponse{
		PostId:                postResult.PostId,
		UploadId:              postResult.PresignedUrl.UploadId,
		PresignedUrls:         postResult.PresignedUrl.ContentPresignedUrls,
		PresignedThumbnailUrl: postResult.PresignedUrl.ThumbanilPresignedUrl,
	}

	api.SendOKWithResult(c, postResponse)
}

func (controller *CreatePostController) ConfirmCreatedPost(c *gin.Context) {
	log.Info().Msg("Handling Request PUT ConfirmCreatedPost")

	var post ConfirmedCreatedPost
	if err := c.BindJSON(&post); err != nil {
		log.Error().Stack().Err(err).Msg("Invalid Data")
		api.SendBadRequest(c, "Invalid Json Request")
		return
	}

	err := controller.service.ConfirmCreatedPost(&post)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOK(c)
}
