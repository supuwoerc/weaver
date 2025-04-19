//go:build wireinject
// +build wireinject

//go:generate wire
package bootstrap

import (
	"gin-web/initialize"
	"gin-web/pkg/email"
	"gin-web/pkg/utils"
	"gin-web/providers"
	"gin-web/router"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.uber.org/zap/zapcore"
	"net/http"
)

func WireApp() *App {
	wire.Build(
		initialize.NewViper,
		initialize.NewWriterSyncer,
		initialize.NewZapLogger,
		initialize.NewDialer,
		initialize.NewEngine,
		initialize.NewServer,
		initialize.NewGORM,
		initialize.NewRedisClient,
		initialize.NewCronLogger,
		initialize.NewCronClient,
		email.NewEmailClient,
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
