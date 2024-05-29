package cmd

import (
	"fmt"
	"gin-web/initialize"
	"gin-web/pkg/global"
)

func Start() {
	initialize.InitConfig()
	global.Logger = initialize.InitZapLogger()
	global.DB = initialize.InitGORM()
	global.RedisClient = initialize.InitRedis()
	handle := initialize.InitEngine()
	initialize.InitServer(handle)
}

func Clean() {
	fmt.Println("关闭服务后的清理...")
}
