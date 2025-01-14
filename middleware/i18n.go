package middleware

import (
	"encoding/json"
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

const (
	CN string = "cn"
	EN string = "en"
)

var languages = []string{CN, EN}

func loadJsonFiles(dir string) ([]*i18n.Message, error) {
	var m []*i18n.Message
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			fileBytes, readErr := os.ReadFile(path)
			if readErr != nil {
				return readErr
			}
			var temp []*i18n.Message
			if readErr = json.Unmarshal(fileBytes, &temp); readErr != nil {
				return readErr
			}
			m = append(m, temp...)
		}
		return nil
	})
	return m, err
}

func I18N() gin.HandlerFunc {
	// 创建一个新的Bundle指定默认语言
	bundle := i18n.NewBundle(language.Chinese)
	// 注册一个JSON加载器
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	// 加载语言文件
	enMessages, err := loadJsonFiles("./pkg/locale/en")
	if err != nil {
		panic(err)
	}
	err = bundle.AddMessages(language.English, enMessages...)
	if err != nil {
		panic(err)
	}
	cnMessages, err := loadJsonFiles("./pkg/locale/zh")
	if err != nil {
		panic(err)
	}
	err = bundle.AddMessages(language.Chinese, cnMessages...)
	if err != nil {
		panic(err)
	}
	defaultLang := viper.GetString("system.defaultLang")
	cnLocalizer := i18n.NewLocalizer(bundle, language.Chinese.String())
	enLocalizer := i18n.NewLocalizer(bundle, language.English.String())
	localeKey := viper.GetString("system.defaultLocaleKey")
	if strings.TrimSpace(localeKey) == "" {
		panic("locale key未配置")
	}
	return func(ctx *gin.Context) {
		locale := ctx.GetHeader(localeKey)
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
		if locale == CN {
			ctx.Set(response.I18nTranslatorKey, cnLocalizer)
		} else if locale == EN {
			ctx.Set(response.I18nTranslatorKey, enLocalizer)
		}
	}
}

func InjectTranslator() gin.HandlerFunc {
	zhTans := zh.New()
	enTans := en.New()
	uni := ut.New(enTans, enTans, zhTans)
	defaultLang := viper.GetString("system.defaultLang")
	zhTrans, _ := uni.GetTranslator("zh")
	enTrans, _ := uni.GetTranslator("en")
	localeKey := viper.GetString("system.defaultLocaleKey")
	if strings.TrimSpace(localeKey) == "" {
		panic("locale key未配置")
	}
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
			locale := ctx.GetHeader(localeKey)
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
			if locale == CN {
				ctx.Set(response.ValidatorTranslatorKey, zhTrans)
			} else if locale == EN {
				ctx.Set(response.ValidatorTranslatorKey, enTrans)
			}
		}
	} else {
		return nil
	}
}
