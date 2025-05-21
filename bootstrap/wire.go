//go:build wireinject
// +build wireinject

//
//go:generate wire
package bootstrap

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	goredislib "github.com/redis/go-redis/v9"
	"github.com/supuwoerc/weaver/initialize"
	"github.com/supuwoerc/weaver/pkg/logger"
	"github.com/supuwoerc/weaver/pkg/utils"
	"github.com/supuwoerc/weaver/providers"
	"github.com/supuwoerc/weaver/router"
	"go.uber.org/zap/zapcore"
	gormLogger "gorm.io/gorm/logger"
)

var loggerProvider = wire.NewSet(
	wire.Bind(new(utils.LocksmithLogger), new(*logger.Logger)),
	wire.Bind(new(initialize.ClientLogger), new(*logger.Logger)),
	logger.NewLogger,
)

var gormLoggerProvider = wire.NewSet(
	wire.Bind(new(gormLogger.Interface), new(*initialize.GormLogger)),
	initialize.NewGormLogger,
)

var redisLoggerProvider = wire.NewSet(
	wire.Bind(new(goredislib.Hook), new(*initialize.RedisLogger)),
	initialize.NewRedisLogger,
)

var emailProvider = wire.NewSet(
	wire.Bind(new(utils.LocksmithEmailClient), new(*initialize.EmailClient)),
	initialize.NewEmailClient,
)

var syncerProvider = wire.NewSet(
	wire.Bind(new(initialize.RedisLogSyncer), new(zapcore.WriteSyncer)),
	wire.Bind(new(initialize.EngineLogger), new(zapcore.WriteSyncer)),
	initialize.NewWriterSyncer,
)

var enginProvider = wire.NewSet(
	wire.Bind(new(http.Handler), new(*gin.Engine)),
	initialize.NewEngine,
)

func WireApp() *App {
	wire.Build(

		initialize.NewViper,
		syncerProvider,
		initialize.NewZapLogger,

		loggerProvider,

		initialize.NewDialer,

		initialize.NewCronLogger,
		initialize.NewCronClient,

		gormLoggerProvider,
		initialize.NewGORM,

		utils.NewRedisLocksmith,

		redisLoggerProvider,
		initialize.NewRedisClient,

		emailProvider,

		enginProvider,
		initialize.NewServer,
		router.NewRouter,

		providers.SystemJobProvider,
		providers.SystemCacheProvider,
		providers.CommonProvider,
		providers.MiddlewareProvider,
		providers.ApiProvider,

		wire.Struct(new(App), "*"),
	)
	return nil
}
