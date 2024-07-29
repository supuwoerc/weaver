package router

import (
	"gin-web/api"
	"gin-web/middleware"
	"github.com/gin-gonic/gin"
)

func InitPingRouter(r *gin.RouterGroup) {
	group := r.Group("ping")
	{
		group.GET("", api.Ping)
		group.GET("exception", api.Exception)
		group.GET("check_permission", middleware.PermissionRequired(), api.CheckPermission)
	}
}
