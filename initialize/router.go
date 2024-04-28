package initialize

import (
	"gin-web/router"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	// 不携带日志和其他中间件
	r := gin.New()
	group := r.Group("api")
	router.InitPingRouter(group)
	router.InitPublicRouter(group)
	return r
}
