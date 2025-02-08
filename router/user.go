package router

import (
	v1 "gin-web/api/v1"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(r *gin.RouterGroup) {
	userApi := v1.NewUserApi()
	userAccessGroup := r.Group("user")
	{
		userAccessGroup.GET("refresh-token")
		userAccessGroup.GET("profile", userApi.Profile)
		userAccessGroup.GET("list", userApi.GetUserList)
		//userAccessGroup.GET("detail", userApi.GetUserDetail)
		//userAccessGroup.POST("update", userApi.UpdateUser)
	}
}
