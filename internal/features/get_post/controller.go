package get_post

import (
	"postservice/internal/api"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type GetPostController struct {
	service *GetPostService
}

type GetPostResponse struct {
	PostUrls []PostUrl `json:"urlPosts"`
}

func NewGetPostController(repository Repository) *GetPostController {
	return &GetPostController{
		service: NewGetPostService(repository),
	}
}

func (controller *GetPostController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/user-posts/:username", controller.GetUserPosts)
}

func (controller *GetPostController) GetUserPosts(c *gin.Context) {
	log.Info().Msg("Handling Request GET UserPosts")
	id := c.Param("username")
	username := string(id)

	postUrls, err := controller.service.GetUserPosts(username)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOKWithResult(c, &GetPostResponse{
		PostUrls: postUrls,
	})
}
