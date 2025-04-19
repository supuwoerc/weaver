package bootstrap

import (
	v1 "gin-web/api/v1"
	"gin-web/conf"
	"gin-web/initialize"
	"gin-web/pkg/cache"
	"gin-web/pkg/constant"
	"gin-web/pkg/job"
	"go.uber.org/zap"
	"sync"
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
	if len(a.conf.System.Hooks.Launch) > 0 {
		for _, item := range a.conf.System.Hooks.Launch {
			switch item {
			case constant.RegisterJobs:
				if err := a.jobManager.RegisterJobsAndStart(); err != nil {
					panic(err)
				}
			case constant.RefreshDeptCache:
				if err := a.cacheManager.Refresh(constant.RefreshDeptCache); err != nil {
					panic(err)
				}
			}
		}
	}
	a.httpServer.Run()
}

func (a *App) Close() {
	defer a.logger.Info("app clean is executed")
	group := sync.WaitGroup{}
	group.Add(1)
	go a.jobManager.Stop(&group)
	group.Wait()
}
