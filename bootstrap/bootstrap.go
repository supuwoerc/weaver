package bootstrap

import (
	"gin-web/initialize"
	"gin-web/pkg/global"
)

func Start() {
	initialize.InitConfig()
	global.Logger = initialize.InitZapLogger()
	global.DB = initialize.InitGORM()
	global.RedisClient = initialize.InitRedis(initialize.LoggerSyncer)
	global.Dialer = initialize.InitDialer()
	global.Cron = initialize.InitCron(global.Logger)
	if err := RegisterJobs(global.Cron, global.Logger); err != nil {
		panic(err)
	}
	initialize.InitServer(initialize.InitEngine(initialize.LoggerSyncer))
}
