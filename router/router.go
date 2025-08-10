package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supuwoerc/weaver/conf"
	"github.com/supuwoerc/weaver/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	swagRoutePattern = "swagger/*any"
)

func NewRouter(r *gin.Engine, conf *conf.Config) *gin.RouterGroup {
	group := r.Group("api/v1")
	// 国际化中间件
	i18n := middleware.NewI18NMiddleware(conf)
	group.Use(i18n.I18N(), i18n.InjectTranslator())
	return group
}

func NotFoundHandler(context *gin.Context) {
	context.HTML(http.StatusNotFound, "404.html", nil)
}

func InitSystemWebRouter(r *gin.Engine) {
	r.NoRoute(NotFoundHandler)
}

func InitSwagWebRouter(r *gin.Engine, conf *conf.Config) {
	if !conf.IsProd() {
		r.GET(swagRoutePattern, ginSwagger.WrapHandler(swaggerFiles.Handler))
	} else {
		r.GET(swagRoutePattern, NotFoundHandler)
	}
}
