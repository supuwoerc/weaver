package router

import "github.com/gin-gonic/gin"

func InitRouter(r *gin.Engine) {
	group := r.Group("api/v1")
	// 系统基础测试
	InitPingRouter(group)
	// swagger文档
	InitSwagger(r)
	// 开放api(不需要走鉴权中间件)
	InitPublicRouter(group)
	// 用户模块
	InitUserRouter(group)
}
