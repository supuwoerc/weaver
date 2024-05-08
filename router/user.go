package router

import (
	"github.com/gin-gonic/gin"
)

func InitUserRouter(r *gin.RouterGroup) {
	//userApi := v1.NewUserApi()
	_ = r.Group("user")
	{

	}
}
