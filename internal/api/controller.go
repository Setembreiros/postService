package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	Routes(routerGroup *gin.RouterGroup)
}

type response struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Content any    `json:"content"`
}

func SendOKWithResult(c *gin.Context, result any) {
	var payload response

	payload.Error = false
	payload.Message = "200 OK"
	payload.Content = result

	c.IndentedJSON(http.StatusOK, payload)
}

func SendFailure(c *gin.Context, httpStatus int, errorMessage string) {
	var payload response

	payload.Error = true
	payload.Message = errorMessage

	c.IndentedJSON(httpStatus, payload)
}

func SendNotFound(c *gin.Context, errorMessage string) {
	SendFailure(c, http.StatusNotFound, errorMessage)
}

func SendInternalServerError(c *gin.Context, errorMessage string) {
	SendFailure(c, http.StatusInternalServerError, errorMessage)
}
