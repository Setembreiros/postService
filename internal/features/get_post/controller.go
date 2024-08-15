package get_post

import (
	"errors"
	"postservice/internal/api"
	database "postservice/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type GetPostController struct {
	service *GetPostService
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

	presignedUrls, err := controller.service.GetUserPosts(username)
	if err != nil {
		var notFoundError *database.NotFoundError
		if errors.As(err, &notFoundError) {
			message := "User not found for username " + username
			api.SendNotFound(c, message)
		} else {
			api.SendInternalServerError(c, err.Error())
		}
		return
	}

	api.SendOKWithResult(c, &presignedUrls)
}
