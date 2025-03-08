package initialize

import (
	"fmt"
	"gin-web/middleware"
	"gin-web/pkg/email"
	"gin-web/router"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io"
	"time"
)

func getEnginLoggerConfig(output io.Writer) gin.LoggerConfig {
	return gin.LoggerConfig{
		Output: output,
		Formatter: func(params gin.LogFormatterParams) string {
			if params.Latency > time.Minute {
				params.Latency = params.Latency.Truncate(time.Second)
			}
			return fmt.Sprintf(
				"{\"caller\":GIN,\"time\":\"%s\",\"status\":%3d,\"latency\":%v,\"method\":\"%s\",\"path\":\"%s\",\"client\":\"%s\"}\n",
				params.TimeStamp.Format(time.DateTime),
				params.StatusCode,
				params.Latency,
				params.Method,
				params.Path,
				params.ClientIP,
			)
		},
	}
}

type EngineLoggerWriter io.Writer

func NewEngine(writer EngineLoggerWriter, emailClient *email.EmailClient, logger *zap.SugaredLogger, v *viper.Viper) *gin.Engine {
	initDebugPrinter(writer)
	// 不携带日志和Recovery中间件，自己添加中间件，为了方便收集Recovery日志
	r := gin.New()
	// html 模板
	r.LoadHTMLGlob(v.GetString("system.templateDir"))
	// 开启ContextWithFallback
	r.ContextWithFallback = true
	// 设置上传文件的最大字节数,Gin默认为32Mb
	maxMultipartMemory := v.GetInt64("system.maxMultipartMemory")
	if maxMultipartMemory > 0 {
		r.MaxMultipartMemory = maxMultipartMemory
	}
	// logger中间件,输出到控制台和zap的日志文件中
	r.Use(gin.LoggerWithConfig(getEnginLoggerConfig(writer)))
	// recovery中间件
	r.Use(middleware.NewRecoveryMiddleware(emailClient, logger, v).Recovery())
	// 跨域中间件
	r.Use(middleware.NewCorsMiddleware(v).Cors())
	// 注册 API 路由
	router.InitApiRouter(r, v)
	// 系统路由
	router.InitSystemWebRouter(r)
	return r
}

func initDebugPrinter(writer io.Writer) {
	gin.DebugPrintFunc = func(format string, values ...interface{}) {
		_, _ = fmt.Fprint(writer, fmt.Sprintf("{\"caller\":GIN DEBUG,\"message\":\"%s\"}\n", fmt.Sprintf(format, values...)))
	}
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		_, _ = fmt.Fprint(writer, fmt.Sprintf("{\"caller\":GIN ROUTER DEBUG,\"method\":\"%s\",\"path\":\"%s\",\"handler\":\"%s\",\"handlers\":%d}\n", httpMethod, absolutePath, handlerName, nuHandlers))
	}
}
