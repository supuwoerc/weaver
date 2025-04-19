package initialize

import (
	"gin-web/conf"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
)

const (
	configType = "yml"
	configPath = "./config"
)

func NewViper() *conf.Config {
	v := viper.New()
	v.SetConfigType(configType)
	v.AddConfigPath(configPath)
	v.SetConfigName("default")
	// 读取默认配置
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	// 设置当前配置环境
	env := determineEnvironment()
	v.SetConfigName(env)
	// 合并配置
	if err := v.MergeInConfig(); err != nil {
		panic(err)
	}
	// 映射到配置对象中
	var globalConfig conf.Config
	if err := v.Unmarshal(&globalConfig); err != nil {
		panic(err)
	}
	return &globalConfig
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
