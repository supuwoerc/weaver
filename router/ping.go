package router

import (
	"gin-web/api"
	"gin-web/middleware"
	"github.com/gin-gonic/gin"
)

func InitPingApiRouter(r *gin.RouterGroup) {
	group := r.Group("ping")
	{
		group.GET("", api.Ping)
		group.GET("exception", api.Exception)
		group.GET("check-permission", middleware.PermissionRequired(), api.CheckPermission)
		group.GET("slow", api.SlowResponse)
		group.GET("check-lock", api.LockResponse)
	}
}
