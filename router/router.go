package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/supuwoerc/weaver/conf"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	swagRoutePattern = "swagger/*any"
)

func NewRouter(r *gin.Engine) *gin.RouterGroup {
	return r.Group("api/v1")
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

func InitPrometheusRouter(r *gin.Engine) {
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

func InitHealthCheckRouter(r *gin.Engine) {
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.DateTime),
		})
	})
}
