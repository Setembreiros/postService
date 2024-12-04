package get_post

import (
	"postservice/internal/api"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type GetPostController struct {
	service *GetPostService
}

type GetPostResponse struct {
	PostUrls   []PostUrl `json:"urlPosts"`
	Limit      int       `json:"limit"`
	NextPostId string    `json:"nextPostId"`
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
	username := c.Param("username")
	lastPostId := c.DefaultQuery("lastPostId", "")
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "6"))

	if err != nil || limit <= 0 {
		api.SendBadRequest(c, "Parámetros de páxinación inválidos")
		return
	}

	postUrls, nextPostId, err := controller.service.GetUserPosts(username, lastPostId, limit)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOKWithResult(c, &GetPostResponse{
		PostUrls:   postUrls,
		Limit:      limit,
		NextPostId: nextPostId,
	})
}
