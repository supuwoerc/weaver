package cmd

import (
	"fmt"
	"gin-web/initialize"
	"gin-web/pkg/global"
)

func Start() {
	global.Logger = initialize.InitZapLogger()
	router := initialize.InitRouter()
	initialize.InitServer(router)
}

func Clean() {
	fmt.Println("关闭服务后的清理...")
}
