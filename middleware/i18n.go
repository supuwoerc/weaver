package middleware

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/samber/lo"
	"github.com/supuwoerc/weaver/conf"
	"github.com/supuwoerc/weaver/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const (
	CN string = "cn"
	EN string = "en"

	defaultZhJsonPath = "./pkg/locale/zh"
	defaultEnJsonPath = "./pkg/locale/en"
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

type I18NMiddleware struct {
	conf *conf.Config
}

func NewI18NMiddleware(conf *conf.Config) *I18NMiddleware {
	return &I18NMiddleware{
		conf: conf,
	}
}

func (i *I18NMiddleware) I18N() gin.HandlerFunc {
	// 创建一个新的Bundle指定默认语言
	bundle := i18n.NewBundle(language.Chinese)
	// 注册一个JSON加载器
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	// 加载语言文件
	enJsonPath := lo.Ternary(i.conf.System.LocalePath.En == "", defaultEnJsonPath, i.conf.System.LocalePath.En)
	enMessages, err := loadJsonFiles(enJsonPath)
	if err != nil {
		panic(err)
	}
	err = bundle.AddMessages(language.English, enMessages...)
	if err != nil {
		panic(err)
	}
	zhJsonPath := lo.Ternary(i.conf.System.LocalePath.Zh == "", defaultZhJsonPath, i.conf.System.LocalePath.Zh)
	cnMessages, err := loadJsonFiles(zhJsonPath)
	if err != nil {
		panic(err)
	}
	err = bundle.AddMessages(language.Chinese, cnMessages...)
	if err != nil {
		panic(err)
	}
	defaultLang := i.conf.System.DefaultLang
	cnLocalizer := i18n.NewLocalizer(bundle, language.Chinese.String())
	enLocalizer := i18n.NewLocalizer(bundle, language.English.String())
	localeKey := i.conf.System.DefaultLocaleKey
	if strings.TrimSpace(localeKey) == "" {
		panic("miss locale key")
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
			ctx.Set(string(response.I18nTranslatorKey), cnLocalizer)
		} else if locale == EN {
			ctx.Set(string(response.I18nTranslatorKey), enLocalizer)
		}
	}
}

func (i *I18NMiddleware) InjectTranslator() gin.HandlerFunc {
	zhTans := zh.New()
	enTans := en.New()
	uni := ut.New(enTans, enTans, zhTans)
	defaultLang := i.conf.System.DefaultLang
	zhTrans, _ := uni.GetTranslator("zh")
	enTrans, _ := uni.GetTranslator("en")
	localeKey := i.conf.System.DefaultLocaleKey
	if strings.TrimSpace(localeKey) == "" {
		panic("missing default locale key configuration")
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
				ctx.Set(string(response.ValidatorTranslatorKey), zhTrans)
			} else if locale == EN {
				ctx.Set(string(response.ValidatorTranslatorKey), enTrans)
			}
		}
	} else {
		return nil
	}
}
