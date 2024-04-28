package initialize

import (
	"gin-web/middleware"
	"gin-web/router"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	// 携带日志和Recovery中间件
	r := gin.Default()
	r.Use(middleware.Cors())
	group := r.Group("api")
	router.InitPingRouter(group)
	router.InitPublicRouter(group)
	return r
}
