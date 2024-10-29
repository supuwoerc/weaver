package bootstrap

import (
	"fmt"
	"gin-web/initialize"
	"gin-web/pkg/global"
)

func Start() {
	initialize.InitConfig()
	logger, writer := initialize.InitZapLogger()
	global.Logger = logger
	global.DB = initialize.InitGORM()
	global.RedisClient = initialize.InitRedis()
	global.Localizer = initialize.InitI18N()
	handle := initialize.InitEngine(writer)
	initialize.InitServer(handle)
}

func Clean() {
	fmt.Println("关闭服务后的清理...")
}
