package initialize

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
)

const (
	_CONFIG_TYPE = "yml"
)

func InitConfig() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configPath := fmt.Sprintf("%s%s", dir, "/config")
	v := viper.New()
	v.SetConfigType(_CONFIG_TYPE)
	v.AddConfigPath(configPath)
	v.SetConfigName("default")
	err = v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	defaultConfig := v.AllSettings()
	for key, val := range defaultConfig {
		viper.SetDefault(key, val)
	}
	env := os.Getenv("GIN_MODE")
	if env == "" || env == gin.DebugMode {
		env = "dev"
	} else if env == gin.TestMode {
		env = "test"
	} else if env == gin.ReleaseMode {
		env = "prod"
	} else {
		panic(errors.New("读取配置文件出错,请检查环境变量:GIN_MODE"))
	}
	viper.SetConfigName(env)
	viper.SetConfigType(_CONFIG_TYPE)
	viper.AddConfigPath(configPath)
	e := viper.ReadInConfig()
	if e != nil {
		panic(e)
	}
}
