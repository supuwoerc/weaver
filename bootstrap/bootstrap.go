package bootstrap

import (
	"context"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/hashicorp/consul/api"
	"github.com/supuwoerc/weaver/api/v1/attachment"
	"github.com/supuwoerc/weaver/api/v1/captcha"
	"github.com/supuwoerc/weaver/api/v1/department"
	"github.com/supuwoerc/weaver/api/v1/permission"
	"github.com/supuwoerc/weaver/api/v1/ping"
	"github.com/supuwoerc/weaver/api/v1/role"
	"github.com/supuwoerc/weaver/api/v1/user"
	"github.com/supuwoerc/weaver/conf"
	"github.com/supuwoerc/weaver/initialize"
	"github.com/supuwoerc/weaver/pkg/cache"
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/consul"
	"github.com/supuwoerc/weaver/pkg/job"
	"github.com/supuwoerc/weaver/pkg/logger"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

const (
	healthCheckEndpoint string = "/health"
)

type App struct {
	logger              *logger.Logger
	conf                *conf.Config
	jobManager          *job.SystemJobManager
	cacheManager        *cache.SystemCacheManager
	elasticsearchClient *elasticsearch.TypedClient
	consulClient        *api.Client
	serviceRegister     *consul.ServiceRegister
	httpServer          *initialize.HttpServer
	traceSpanExporter   tracesdk.SpanExporter
	tracerProvider      *tracesdk.TracerProvider
	attachmentApi       *attachment.Api
	captchaApi          *captcha.Api
	departmentApi       *department.Api
	permissionApi       *permission.Api
	pingApi             *ping.Api
	roleApi             *role.Api
	userApi             *user.Api
}

func (a *App) Run() {
	// 注册定时任务
	if err := a.jobManager.RegisterJobsAndStart(); err != nil {
		panic(err)
	}
	// 执行相关hook
	if len(a.conf.System.Hooks.Launch) > 0 {
		for _, item := range a.conf.System.Hooks.Launch {
			switch item {
			case constant.AutoManageDeptCache:
				if err := a.cacheManager.Refresh(context.Background(), constant.AutoManageDeptCache); err != nil {
					panic(err)
				}
			}
		}
	}
	// 注册服务
	if err := a.registerService(); err != nil {
		panic(err)
	}
	a.httpServer.Run()
}

func (a *App) Close() {
	defer func() {
		// 关闭 oltp exporter
		err := a.traceSpanExporter.Shutdown(context.Background())
		if err != nil {
			a.logger.Errorw("Error shutting down tracer provider", "err", err.Error())
		}
		// 日志相关 sync
		_ = a.logger.Sync()
	}()
	defer a.logger.Info("app clean is executed")
	// 先下线服务,避免外部继续调用服务
	if err := a.serviceRegister.Deregister(); err != nil {
		a.logger.Errorw("Failed to deregister service", "err", err.Error())
	}
	// 停止定时任务
	a.jobManager.Stop()
	// 执行相关 hook
	if len(a.conf.System.Hooks.Close) > 0 {
		for _, item := range a.conf.System.Hooks.Close {
			switch item {
			case constant.AutoManageDeptCache:
				if err := a.cacheManager.Clean(context.Background(), constant.AutoManageDeptCache); err != nil {
					panic(err)
				}
			}
		}
	}
}

// 服务注册
func (a *App) registerService() error {
	interval := 15 * time.Second
	if a.conf.System.HealthCheckInterval > 0 {
		interval = a.conf.System.HealthCheckInterval * time.Second
	}
	return a.serviceRegister.Register(a.conf.AppName, a.httpServer.Addr(), a.httpServer.Port(),
		consul.WithTags(a.conf.AppName, a.conf.AppVersion, a.conf.Env),
		consul.WithMeta(map[string]string{
			"version": a.conf.AppVersion,
			"env":     a.conf.Env,
		}),
		consul.WithHTTPCheck(healthCheckEndpoint, interval),
	)
}

type Cli struct {
	Logger *logger.Logger
	Conf   *conf.Config
}
