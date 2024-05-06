package router

import (
	"gin-web/api"
	"github.com/gin-gonic/gin"
)

func InitPingRouter(r *gin.RouterGroup) *gin.RouterGroup {
	group := r.Group("ping")
	{
		group.GET("", api.Ping)
		group.GET("exception", api.Exception)
	}
	return group
}
