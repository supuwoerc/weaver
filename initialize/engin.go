package initialize

import (
	"github.com/supuwoerc/weaver/conf"
	"github.com/supuwoerc/weaver/middleware"
	"github.com/supuwoerc/weaver/pkg/logger"
	local "github.com/supuwoerc/weaver/pkg/redis"
	"github.com/supuwoerc/weaver/router"

	"github.com/gin-gonic/gin"
)

func NewEngine(emailClient *EmailClient, rc *local.CommonRedisClient, logger *logger.Logger, conf *conf.Config) *gin.Engine {
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
	// 国际化中间件
	i18n := middleware.NewI18NMiddleware(conf)
	r.Use(i18n.I18N(), i18n.InjectTranslator())
	// trace中间件,在上下文中放入trace信息
	r.Use(middleware.NewTraceMiddleware(conf, logger).Trace())
	// logger中间件,输出到控制台和zap的日志文件中
	r.Use(middleware.NewEnginLoggerMiddleware(logger).Logger())
	// recovery中间件
	r.Use(middleware.NewRecoveryMiddleware(emailClient, logger, conf).Recovery())
	// 跨域中间件
	r.Use(middleware.NewCorsMiddleware(conf).Cors())
	// prometheus监控
	r.Use(middleware.NewPrometheusMiddleware().Prometheus())
	// 开启ForwardedByClientIP(配合限流)
	r.ForwardedByClientIP = true
	// 系统限流中间件
	r.Use(middleware.NewLimiterMiddleware(rc, conf).RequestLimit())
	// 系统相关路由
	router.InitSystemWebRouter(r)
	// swag相关路由
	router.InitSwagWebRouter(r, conf)
	// prometheus相关路由
	router.InitPrometheusRouter(r)
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
