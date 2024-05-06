package router

import "github.com/gin-gonic/gin"

func InitRouter(r *gin.Engine) {
	group := r.Group("api/v1")
	InitPingRouter(group)
	InitPublicRouter(group)
	InitUserRouter(group)
}
