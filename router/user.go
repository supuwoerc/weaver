package router

import (
	v1 "gin-web/api/v1"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(r *gin.RouterGroup) {
	userApi := v1.NewUserApi()
	userAccessGroup := r.Group("user")
	{
		userAccessGroup.GET("refresh_token")
		userAccessGroup.GET("profile", userApi.Profile)
		// TODO:限制管理员权限
		userAccessGroup.POST("set_roles", userApi.SetRoles)
		// TODO:限制管理员权限
		userAccessGroup.GET("get_roles", userApi.GetRoles)
	}
}
