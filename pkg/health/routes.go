package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Routes struct{}

func NewHealthCheckRoutes() *Routes {
	return &Routes{}
}

func (r *Routes) RegisterRoutes(health *gin.RouterGroup, group *gin.RouterGroup) {
	health.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{"message": "pong"})
	})

	group.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{"message": "pong"})
	})
}
