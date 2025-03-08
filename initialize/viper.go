package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
)

const (
	configType = "yml"
	configPath = "./config"
)

func NewViper() *viper.Viper {
	v := viper.New()
	v.SetConfigType(configType)
	v.AddConfigPath(configPath)
	v.SetConfigName("default")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	defaultConfig := v.AllSettings()
	for key, val := range defaultConfig {
		viper.SetDefault(key, val)
	}
	env := determineEnvironment()
	viper.SetConfigName(env)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)
	e := viper.ReadInConfig()
	if e != nil {
		panic(e)
	}
	return v
}

// 辅助函数判断环境
func determineEnvironment() string {
	switch mode := os.Getenv("GIN_MODE"); mode {
	case gin.ReleaseMode:
		return "prod"
	case gin.TestMode:
		return "test"
	default: // 包含空字符串和 DebugMode
		return "dev"
	}
}
