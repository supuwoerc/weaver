package router

import (
	"gin-web/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitApiRouter(r *gin.Engine) *gin.RouterGroup {
	group := r.Group("api/v1")
	// 国际化中间件
	group.Use(middleware.I18N(), middleware.InjectTranslator())
	return group
}

func InitSystemWebRouter(r *gin.Engine) {
	r.NoRoute(func(context *gin.Context) {
		context.HTML(http.StatusOK, "404.html", nil)
	})
}
