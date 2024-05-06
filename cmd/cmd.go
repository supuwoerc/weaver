package cmd

import (
	"fmt"
	"gin-web/initialize"
	"gin-web/pkg/global"
)

func Start() {
	global.Logger = initialize.InitZapLogger()
	global.DB = initialize.InitGORM()
	handle := initialize.InitEngine()
	initialize.InitServer(handle)
}

func Clean() {
	fmt.Println("关闭服务后的清理...")
}
