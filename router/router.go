package router

import (
	"gin-web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func NewRouter(r *gin.Engine, v *viper.Viper) *gin.RouterGroup {
	group := r.Group("api/v1")
	// 国际化中间件
	i18n := middleware.NewI18NMiddleware(v)
	group.Use(i18n.I18N(), i18n.InjectTranslator())
	return group
}

func InitSystemWebRouter(r *gin.Engine) {
	r.NoRoute(func(context *gin.Context) {
		context.HTML(http.StatusNotFound, "404.html", nil)
	})
}
