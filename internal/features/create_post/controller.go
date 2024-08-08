package create_post

import (
	"fmt"
	"postservice/internal/api"

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

func NewCreatePostController(repository Repository) *CreatePostController {
	return &CreatePostController{
		service: NewCreatePostService(repository),
	}
}

func (controller *CreatePostController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/post", controller.CreatePost)
}

func (controller *CreatePostController) CreatePost(c *gin.Context) {
	log.Info().Msg("Handling Request POST CreatePost")
	var post Post

	if err := c.BindJSON(&post); err != nil {
		log.Error().Stack().Err(err).Msg("Invalid Data")
		return
	}

	postId, presignedUrl, err := controller.service.CreatePost(&post)
	fmt.Println(presignedUrl)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOKWithResult(c, &CreatePostResponse{
		PostId:       postId,
		PresignedUrl: presignedUrl,
	})
}
