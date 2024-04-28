package cmd

import (
	"fmt"
	"gin-web/initialize"
)

func Start() {
	router := initialize.InitRouter()
	initialize.InitServer(router)
}

func Clean() {
	fmt.Println("关闭服务后的清理...")
}
