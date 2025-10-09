//go:build wireinject
// +build wireinject

//go:generate wire
package bootstrap

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	goredislib "github.com/redis/go-redis/v9"
	"github.com/supuwoerc/weaver/initialize"
	"github.com/supuwoerc/weaver/pkg/logger"
	"github.com/supuwoerc/weaver/pkg/utils"
	"github.com/supuwoerc/weaver/providers"
	"github.com/supuwoerc/weaver/router"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	gormLogger "gorm.io/gorm/logger"
)

func WireApp() *App {
	wire.Build(

		initialize.NewViper,

		initialize.NewWriterSyncer,

		initialize.NewZapLogger,

		wire.Bind(new(utils.LocksmithLogger), new(*logger.Logger)),
		wire.Bind(new(initialize.ClientLogger), new(*logger.Logger)),
		logger.NewLogger,

		initialize.NewDialer,

		initialize.NewOLTPExporter,
		wire.Bind(new(tracesdk.SpanExporter), new(*otlptrace.Exporter)),
		initialize.NewTracerProvider,

		initialize.NewCronLogger,
		initialize.NewCronClient,

		wire.Bind(new(gormLogger.Interface), new(*initialize.GormLogger)),
		initialize.NewGormLogger,
		initialize.NewGORM,

		wire.Bind(new(elastictransport.Logger), new(*initialize.ElasticsearchLogger)),
		initialize.NewElasticsearchLogger,
		initialize.NewElasticsearchClient,

		utils.NewRedisLocksmith,

		wire.Bind(new(goredislib.Hook), new(*initialize.RedisLogger)),
		initialize.NewRedisLogger,

		initialize.NewRedisClient,

		initialize.NewEmailClient,

		wire.Bind(new(initialize.OSSClient), new(*s3.Client)),
		initialize.NewS3Client,

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

func WireCli() *Cli {
	wire.Build(
		initialize.NewViper,
		initialize.NewWriterSyncer,

		initialize.NewZapLogger,
		logger.NewLogger,

		wire.Struct(new(Cli), "*"),
	)
	return nil
}
