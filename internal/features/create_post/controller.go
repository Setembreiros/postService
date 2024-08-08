package create_post

import (
	"postservice/internal/api"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CreatePostController struct {
}

func NewCreatePostController() *CreatePostController {
	return &CreatePostController{}
}

func (controller *CreatePostController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/post", controller.CreatePost)
}

func (controller *CreatePostController) CreatePost(c *gin.Context) {
	log.Info().Msg("Handling Request POST CreatePost")

	api.SendOKWithResult(c, "Everything Ok")
}
