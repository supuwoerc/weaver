//go:build wireinject
// +build wireinject

//
//go:generate wire
package bootstrap

import (
	"net/http"

	initialize "github.com/supuwoerc/weaver/initialize"
	"github.com/supuwoerc/weaver/pkg/utils"
	providers "github.com/supuwoerc/weaver/providers"
	router "github.com/supuwoerc/weaver/router"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.uber.org/zap/zapcore"
)

func WireApp() *App {
	wire.Build(
		initialize.NewViper,
		initialize.NewZapLogger,
		initialize.NewWriterSyncer,
		initialize.NewDialer,
		initialize.NewEngine,
		initialize.NewServer,
		initialize.NewGORM,
		initialize.NewRedisClient,
		initialize.NewCronLogger,
		initialize.NewCronClient,
		wire.Bind(new(http.Handler), new(*gin.Engine)),
		wire.Bind(new(initialize.EngineLogger), new(zapcore.WriteSyncer)),
		wire.Bind(new(initialize.RedisLogSyncer), new(zapcore.WriteSyncer)),
		utils.NewRedisLocksmith,
		router.NewRouter,
		providers.CommonProvider,
		providers.SystemCacheProvider,
		providers.SystemJobProvider,
		providers.V1Provider,
		wire.Struct(new(App), "*"),
	)
	return nil
}
