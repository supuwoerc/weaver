//go:build wireinject
// +build wireinject

//go:generate wire
package bootstrap

import (
	"gin-web/initialize"
	"gin-web/pkg/email"
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
		wire.Bind(new(initialize.EngineLoggerWriter), new(zapcore.WriteSyncer)),
		email.NewEmailClient,
		initialize.NewEngine,
		initialize.NewServer,
		wire.Struct(new(App), "*"),
		//api.ApiProvider,
	)
	return nil
}
