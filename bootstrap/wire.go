//go:build wireinject
// +build wireinject

//go:generate wire
package bootstrap

import (
	"gin-web/api"
	"gin-web/initialize"
	"gin-web/pkg/email"
	"gin-web/pkg/utils"
	"gin-web/router"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.uber.org/zap/zapcore"
	"net/http"
)

func wireApp() *App {
	wire.Build(
		initialize.NewViper,
		initialize.NewWriterSyncer,
		initialize.NewZapLogger,
		initialize.NewDialer,
		wire.Bind(new(http.Handler), new(*gin.Engine)),
		wire.Bind(new(initialize.EngineLogger), new(zapcore.WriteSyncer)),
		wire.Bind(new(initialize.RedisClientLogger), new(zapcore.WriteSyncer)),
		email.NewEmailClient,
		initialize.NewEngine,
		initialize.NewServer,
		initialize.NewGORM,
		initialize.NewRedisClient,
		utils.NewRedisLocksmith,
		router.NewApiRouter,
		api.ApiProvider,
		wire.Struct(new(App), "*"),
	)
	return nil
}
