package create_post

import (
	"postservice/internal/api"
	"postservice/internal/bus"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CreatePostController struct {
	service *CreatePostService
}

type CreatePostResponse struct {
	PostId       string `json:"post_id"`
	PresignedUrl string `json:"presigned_url"`
}

func NewCreatePostController(repository Repository, bus *bus.EventBus) *CreatePostController {
	return &CreatePostController{
		service: NewCreatePostService(repository, bus),
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

	postId, presignedUrl, err := controller.service.CreatePost(&post)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOKWithResult(c, &CreatePostResponse{
		PostId:       postId,
		PresignedUrl: presignedUrl,
	})
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
