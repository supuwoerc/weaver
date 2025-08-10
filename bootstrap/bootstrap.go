package bootstrap

import (
	"context"

	v1 "github.com/supuwoerc/weaver/api/v1/attachment"
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
	"github.com/supuwoerc/weaver/pkg/job"
	"github.com/supuwoerc/weaver/pkg/logger"
)

type App struct {
	logger        *logger.Logger
	conf          *conf.Config
	jobManager    *job.SystemJobManager
	cacheManager  *cache.SystemCacheManager
	httpServer    *initialize.HttpServer
	attachmentApi *v1.Api
	captchaApi    *captcha.Api
	departmentApi *department.Api
	permissionApi *permission.Api
	pingApi       *ping.Api
	roleApi       *role.Api
	userApi       *user.Api
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
	a.httpServer.Run()
}

func (a *App) Close() {
	// logger sync
	defer func() {
		_ = a.logger.Sync()
	}()
	defer a.logger.Info("app clean is executed")
	// 停止定时任务
	a.jobManager.Stop()
	// 执行相关hook
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
