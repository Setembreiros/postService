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
	routerGroup.DELETE("/posts", controller.DeletePosts)
}

func (controller *DeletePostController) DeletePosts(c *gin.Context) {
	log.Info().Msg("Handling Request Delete Posts")
	postIds := c.QueryArray("post_id")
	if len(postIds) == 0 {
		api.SendBadRequest(c, "Missing post_id parameters")
		return
	}

	err := controller.service.DeletePosts(postIds)
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
