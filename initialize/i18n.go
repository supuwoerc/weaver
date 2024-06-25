package initialize

import (
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func InitI18N() map[string]*i18n.Localizer {
	// 创建一个新的Bundle指定默认语言
	bundle := i18n.NewBundle(language.SimplifiedChinese)
	// 注册一个JSON加载器
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	// 加载语言文件
	bundle.MustLoadMessageFile("./pkg/locale/en.json")
	bundle.MustLoadMessageFile("./pkg/locale/zh-Hans.json")
	return map[string]*i18n.Localizer{
		"cn": i18n.NewLocalizer(bundle, "zh"),
		"en": i18n.NewLocalizer(bundle, "en"),
	}
}
