package cmd

import (
	"fmt"
	"gin-web/initialize"
	"gin-web/pkg/constant"
	"gin-web/pkg/global"
)

func Start() {
	initialize.InitConfig()
	global.Logger = initialize.InitZapLogger()
	global.DB = initialize.InitGORM()
	global.RedisClient = initialize.InitRedis()
	global.Localizer = initialize.InitI18N()
	global.LocaleErrors = constant.InitErrors(constant.InitParams{
		CN: global.Localizer[global.CN],
		EN: global.Localizer[global.EN],
	})
	global.CasbinEnforcer = initialize.InitCasbin(global.DB)
	handle := initialize.InitEngine()
	initialize.InitServer(handle)
}

func Clean() {
	fmt.Println("关闭服务后的清理...")
}
