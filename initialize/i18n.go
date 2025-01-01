package initialize

import (
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

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

func InitI18N() map[string]*i18n.Localizer {
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
	return map[string]*i18n.Localizer{
		"cn": i18n.NewLocalizer(bundle, language.Chinese.String()),
		"en": i18n.NewLocalizer(bundle, language.English.String()),
	}
}
