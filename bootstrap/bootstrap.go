package bootstrap

import (
	v1 "gin-web/api/v1"
	"gin-web/initialize"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type App struct {
	logger        *zap.SugaredLogger
	viper         *viper.Viper
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
	a.httpServer.Run()
}

func (a *App) Close() {
}

//func Clean() {
//	defer global.Logger.Info("Clean is executed")
//	group := sync.WaitGroup{}
//	group.Add(1)
//	go cleanCronJob(&group)
//	group.Wait()
//}

//func cleanCronJob(group *sync.WaitGroup) {
//	defer group.Done()
//	ctx := global.Cron.Stop()
//	<-ctx.Done()
//	global.Logger.Info("Cron jobs have been stopped")
//}
