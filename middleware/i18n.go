package middleware

import (
	"gin-web/pkg/global"
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var languages = []string{"cn", "en"}

func I18N() gin.HandlerFunc {
	return func(context *gin.Context) {
		locale := context.GetHeader("Locale")
		exist := false
		for _, val := range languages {
			if val == locale {
				exist = true
				break
			}
		}
		if !exist {
			locale = viper.GetString("system.defaultLang")
		}
		context.Set(response.TranslatorKey, global.Localizer[locale])
	}
}
