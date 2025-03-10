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

func RunApp() {
	app := wireApp()
	app.httpServer.Run()
}
func CleanApp() {}

//func Start() {
//	initialize.NewViper()
//	global.Logger = initialize.NewZapLogger()
//	global.DB = initialize.InitGORM()
//	global.RedisClient = initialize.NewRedisClient(initialize.LoggerSyncer)
//	global.Dialer = initialize.InitDialer()
//	global.Cron, global.CronLogger = initialize.InitCron(global.Logger)
//	if err := RegisterJobs(); err != nil {
//		panic(err)
//	}
//	initialize.NewServer(initialize.NewEngine(initialize.LoggerSyncer))
//}
//
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
