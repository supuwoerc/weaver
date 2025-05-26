package initialize

import (
	"github.com/supuwoerc/weaver/conf"
	"github.com/supuwoerc/weaver/middleware"
	"github.com/supuwoerc/weaver/pkg/logger"
	"github.com/supuwoerc/weaver/router"

	"github.com/gin-gonic/gin"
)

func NewEngine(emailClient *EmailClient, logger *logger.Logger, conf *conf.Config) *gin.Engine {
	initDebugLogger(logger)
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
	r.Use(middleware.NewEnginLoggerMiddleware(logger).Logger())
	// recovery中间件
	r.Use(middleware.NewRecoveryMiddleware(emailClient, logger, conf).Recovery())
	// 跨域中间件
	r.Use(middleware.NewCorsMiddleware(conf).Cors())
	// 系统相关路由
	router.InitSystemWebRouter(r)
	return r
}

func initDebugLogger(logger *logger.Logger) {
	gin.DebugPrintFunc = func(format string, values ...interface{}) {
		logger.Infof(format, values...)
	}
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		logger.Infow("route register",
			"method", httpMethod,
			"path", absolutePath,
			"name", handlerName,
			"handlers", nuHandlers,
		)
	}
}
