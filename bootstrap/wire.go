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

func WireApp() *App {
	wire.Build(

		initialize.NewViper,

		wire.Bind(new(initialize.RedisLogSyncer), new(zapcore.WriteSyncer)),
		wire.Bind(new(initialize.EngineLogger), new(zapcore.WriteSyncer)),
		initialize.NewWriterSyncer,

		initialize.NewZapLogger,

		wire.Bind(new(utils.LocksmithLogger), new(*logger.Logger)),
		wire.Bind(new(initialize.ClientLogger), new(*logger.Logger)),
		logger.NewLogger,

		initialize.NewDialer,

		initialize.NewCronLogger,
		initialize.NewCronClient,

		wire.Bind(new(gormLogger.Interface), new(*initialize.GormLogger)),
		initialize.NewGormLogger,
		initialize.NewGORM,

		utils.NewRedisLocksmith,

		wire.Bind(new(goredislib.Hook), new(*initialize.RedisLogger)),
		initialize.NewRedisLogger,

		initialize.NewRedisClient,

		wire.Bind(new(utils.LocksmithEmailClient), new(*initialize.EmailClient)),
		initialize.NewEmailClient,

		wire.Bind(new(http.Handler), new(*gin.Engine)),
		initialize.NewEngine,
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
