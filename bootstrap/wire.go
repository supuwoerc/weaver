//go:build wireinject
// +build wireinject

//
//go:generate wire
package bootstrap

import (
	"net/http"

	initialize "github.com/supuwoerc/weaverinitialize"
	"github.com/supuwoerc/weaverpkg/utils"
	providers "github.com/supuwoerc/weaverproviders"
	router "github.com/supuwoerc/weaverrouter"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.uber.org/zap/zapcore"
)

func WireApp() *App {
	wire.Build(
		initialize.NewViper,
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
		wire.Bind(new(initialize.RedisClientLogger), new(zapcore.WriteSyncer)),
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
