package initialize

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/supuwoerc/weaver/conf"
	"github.com/supuwoerc/weaver/middleware"
	"github.com/supuwoerc/weaver/pkg/logger"
	"github.com/supuwoerc/weaver/router"

	"github.com/gin-gonic/gin"
)

func getEnginLoggerConfig(output io.Writer) gin.LoggerConfig {
	return gin.LoggerConfig{
		Output: output,
		Formatter: func(params gin.LogFormatterParams) string {
			if params.Latency > time.Minute {
				params.Latency = params.Latency.Truncate(time.Second)
			}
			var builder strings.Builder
			builder.WriteString(`{"caller":"GIN","time":"`)
			builder.WriteString(params.TimeStamp.Format(time.DateTime))
			builder.WriteString(`","status":`)
			builder.WriteString(strconv.Itoa(params.StatusCode))
			builder.WriteString(`,"latency":`)
			builder.WriteString(params.Latency.String())
			builder.WriteString(`,"method":"`)
			builder.WriteString(params.Method)
			builder.WriteString(`","path":"`)
			builder.WriteString(params.Path)
			if traceId, ok := params.Keys[string(logger.TraceIdContextKey)]; ok {
				tid, o := traceId.(string)
				if o {
					builder.WriteString(`","trace_id":"`)
					builder.WriteString(tid)
				}
			}
			builder.WriteString(`","client":"`)
			builder.WriteString(params.ClientIP)
			builder.WriteString(`"}`)
			builder.WriteByte('\n') // 换行符
			return builder.String()
		},
	}
}

type EngineLogger io.Writer

func NewEngine(writer EngineLogger, emailClient *EmailClient, logger *logger.Logger, conf *conf.Config) *gin.Engine {
	initDebugPrinter(writer)
	// 不携带日志和Recovery中间件，自己添加中间件，为了方便收集Recovery日志
	r := gin.New()
	// html 模板
	r.LoadHTMLGlob(conf.System.TemplateDir)
	// 开启ContextWithFallback
	r.ContextWithFallback = true
	// 设置上传文件的最大字节数,Gin默认为32Mb
	maxMultipartMemory := conf.System.MaxMultipartMemory
	if maxMultipartMemory > 0 {
		r.MaxMultipartMemory = maxMultipartMemory
	}
	// trace中间件,在上下文中放入trace信息
	r.Use(middleware.NewTraceMiddleware(conf, logger).Trace())
	// logger中间件,输出到控制台和zap的日志文件中
	r.Use(gin.LoggerWithConfig(getEnginLoggerConfig(writer)))
	// recovery中间件
	r.Use(middleware.NewRecoveryMiddleware(emailClient, logger, conf).Recovery())
	// 跨域中间件
	r.Use(middleware.NewCorsMiddleware(conf).Cors())
	// 系统相关路由
	router.InitSystemWebRouter(r)
	return r
}

func initDebugPrinter(writer io.Writer) {
	gin.DebugPrintFunc = func(format string, values ...interface{}) {
		_, _ = fmt.Fprintf(writer, "{\"caller\":GIN DEBUG,\"message\":\"%s\"}\n", fmt.Sprintf(format, values...))
	}
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		_, _ = fmt.Fprintf(writer, "{\"caller\":GIN ROUTER DEBUG,\"method\":\"%s\",\"path\":\"%s\",\"handler\":\"%s\",\"handlers\":%d}\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
}
