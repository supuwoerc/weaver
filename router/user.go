package router

import (
	v1 "gin-web/api/v1"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(r *gin.RouterGroup) *gin.RouterGroup {
	userApi := v1.NewUserApi()
	group := r.Group("user")
	{
		group.POST("signup", userApi.SignUp)
	}
	return group
}
