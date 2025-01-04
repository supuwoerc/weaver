package initialize

import (
	"gin-web/middleware"
	"gin-web/router"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
)

func InitEngine(writer io.Writer) *gin.Engine {
	// 不携带日志和Recovery中间件，自己添加中间件，为了方便收集Recovery日志
	r := gin.New()
	// 开启ContextWithFallback
	r.ContextWithFallback = true
	// 设置上传文件的最大字节数,Gin默认为32Mb
	maxMultipartMemory := viper.GetInt64("system.maxMultipartMemory")
	if maxMultipartMemory > 0 {
		r.MaxMultipartMemory = maxMultipartMemory
	}
	// logger中间件,输出到控制台和zap的日志文件中
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Output: writer,
	}))
	// recovery中间件
	r.Use(middleware.Recovery())
	// 跨域中间件
	r.Use(middleware.Cors())
	// 注册路由
	router.InitRouter(r)
	return r
}
