package bootstrap

import (
	v1 "gin-web/api/v1"
	"gin-web/initialize"
	"gin-web/pkg/job"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"sync"
)

type App struct {
	logger        *zap.SugaredLogger
	viper         *viper.Viper
	jobRegisterer *job.JobRegisterer
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
	if err := a.jobRegisterer.RegisterJobsAndStart(); err != nil {
		panic(err)
	}
	a.httpServer.Run()
}

func (a *App) Close() {
	defer a.logger.Info("app clean is executed")
	group := sync.WaitGroup{}
	group.Add(1)
	go a.jobRegisterer.Stop(&group)
	group.Wait()
}
