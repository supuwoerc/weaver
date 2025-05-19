package bootstrap

import (
	"context"

	v1 "github.com/supuwoerc/weaver/api/v1"
	"github.com/supuwoerc/weaver/conf"
	"github.com/supuwoerc/weaver/initialize"
	"github.com/supuwoerc/weaver/pkg/cache"
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/job"

	"go.uber.org/zap"
)

type App struct {
	logger        *zap.SugaredLogger
	conf          *conf.Config
	jobManager    *job.SystemJobManager
	cacheManager  *cache.SystemCacheManager
	httpServer    *initialize.HttpServer
	attachmentApi *v1.AttachmentApi
	captchaApi    *v1.CaptchaApi
	departmentApi *v1.DepartmentApi
	permissionApi *v1.PermissionApi
	pingApi       *v1.PingApi
	roleApi       *v1.RoleApi
	userApi       *v1.UserApi
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
