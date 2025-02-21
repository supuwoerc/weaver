package bootstrap

import (
	"gin-web/initialize"
	"gin-web/pkg/global"
	"sync"
)

func Start() {
	initialize.InitConfig()
	global.Logger = initialize.InitZapLogger()
	global.DB = initialize.InitGORM()
	global.RedisClient = initialize.InitRedis(initialize.LoggerSyncer)
	global.Dialer = initialize.InitDialer()
	global.Cron, global.CronLogger = initialize.InitCron(global.Logger)
	if err := RegisterJobs(); err != nil {
		panic(err)
	}
	initialize.InitServer(initialize.InitEngine(initialize.LoggerSyncer))
}

func Clean() {
	group := sync.WaitGroup{}
	group.Add(1)
	go cleanCronJob(&group)
	group.Wait()
}

func cleanCronJob(group *sync.WaitGroup) {
	defer group.Done()
	ctx := global.Cron.Stop()
	<-ctx.Done()
	global.Logger.Info("Cron jobs have been stopped")
}
