package middleware

import (
	"gin-web/pkg/global"
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var languages = []string{global.CN, global.EN}

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
		context.Set(response.Locale, locale)
	}
}
