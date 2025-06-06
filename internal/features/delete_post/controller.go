package delete_post

import (
	"errors"
	"fmt"
	"postservice/internal/api"
	"postservice/internal/bus"
	database "postservice/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type DeletePostController struct {
	service *DeletePostService
}

func NewDeletePostController(repository Repository, bus *bus.EventBus) *DeletePostController {
	return &DeletePostController{
		service: NewDeletePostService(repository, bus),
	}
}

func (controller *DeletePostController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.DELETE("/posts/:username", controller.DeletePosts)
}

func (controller *DeletePostController) DeletePosts(c *gin.Context) {
	log.Info().Msg("Handling Request Delete Posts")
	username := c.Param("username")
	postIds := c.QueryArray("postId")
	if len(postIds) == 0 {
		api.SendBadRequest(c, "Missing postId parameters")
		return
	}

	err := controller.service.DeletePosts(username, postIds)
	if err != nil {
		var notFoundError *database.NotFoundError
		if errors.As(err, &notFoundError) {
			message := fmt.Sprintf("Some posts were not found for post ids %v", postIds)
			api.SendNotFound(c, message)
		} else {
			api.SendInternalServerError(c, err.Error())
		}
		return
	}

	api.SendOK(c)
}
