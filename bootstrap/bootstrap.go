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
	global.Localizer = initialize.InitI18N()
	global.Dialer = initialize.InitDialer()
	initialize.InitServer(initialize.InitEngine(initialize.LoggerSyncer))
}
