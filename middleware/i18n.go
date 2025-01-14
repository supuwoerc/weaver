package middleware

import (
	"gin-web/pkg/global"
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/spf13/viper"
	"reflect"
	"strings"
)

var languages = []string{global.CN, global.EN}

func I18N() gin.HandlerFunc {
	defaultLang := viper.GetString("system.defaultLang")
	return func(ctx *gin.Context) {
		locale := ctx.GetHeader("Locale")
		exist := false
		for _, val := range languages {
			if val == locale {
				exist = true
				break
			}
		}
		if !exist {
			locale = defaultLang
		}
		// 在上下文中注入i18n.Localizer实例
		ctx.Set(response.I18nTranslatorKey, global.Localizer[locale])
	}
}

func InjectTranslator() gin.HandlerFunc {
	zhTans := zh.New()
	enTans := en.New()
	uni := ut.New(enTans, enTans, zhTans)
	defaultLang := viper.GetString("system.defaultLang")
	zhTrans, _ := uni.GetTranslator("zh")
	enTrans, _ := uni.GetTranslator("en")
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			if name != "" {
				return name
			} else {
				name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
				return name
			}
		})
		if err := zhTranslations.RegisterDefaultTranslations(v, zhTrans); err != nil {
			panic(err)
		}
		if err := enTranslations.RegisterDefaultTranslations(v, enTrans); err != nil {
			panic(err)
		}
		return func(ctx *gin.Context) {
			locale := ctx.GetHeader("Locale")
			exist := false
			for _, val := range languages {
				if val == locale {
					exist = true
					break
				}
			}
			if !exist {
				locale = defaultLang
			}
			if locale == global.CN {
				ctx.Set(response.ValidatorTranslatorKey, zhTrans)
			} else if locale == global.EN {
				ctx.Set(response.ValidatorTranslatorKey, enTrans)
			}
		}
	} else {
		return nil
	}
}
