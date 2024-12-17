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
	PostUrls          []PostUrl `json:"urlPosts"`
	Limit             int       `json:"limit"`
	LastPostId        string    `json:"lastPostId"`
	LastPostCreatedAt string    `json:"lastPostCreatedAt"`
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
	lastPostCreatedAt := c.DefaultQuery("lastPostCreatedAt", "")
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "6"))

	if err != nil || limit <= 0 {
		api.SendBadRequest(c, "Invalid pagination parameters, limit has to be greater than 0")
		return
	}

	if (lastPostId != "" && lastPostCreatedAt == "") || (lastPostId == "" && lastPostCreatedAt != "") {
		api.SendBadRequest(c, "Invalid pagination parameters, lastPostId and lastPostCreatedAt both have to have value or both have to be empty")
		return
	}

	postUrls, lastPostId, lastPostCreatedAt, err := controller.service.GetUserPosts(username, lastPostId, lastPostCreatedAt, limit)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOKWithResult(c, &GetPostResponse{
		PostUrls:          postUrls,
		Limit:             limit,
		LastPostId:        lastPostId,
		LastPostCreatedAt: lastPostCreatedAt,
	})
}
