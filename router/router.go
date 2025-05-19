package router

import (
	"gin-web/conf"
	"gin-web/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine, conf *conf.Config) *gin.RouterGroup {
	group := r.Group("api/v1")
	// 国际化中间件
	i18n := middleware.NewI18NMiddleware(conf)
	group.Use(i18n.I18N(), i18n.InjectTranslator())
	return group
}

func InitSystemWebRouter(r *gin.Engine) {
	r.NoRoute(func(context *gin.Context) {
		context.HTML(http.StatusNotFound, "404.html", nil)
	})
}
