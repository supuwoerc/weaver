package main

import (
	"gin-web/bootstrap"
	"gin-web/cmd"
)

var isCli = false

// @title          Learn-Gin-Web API
// @version        1.0
// @description    This is a sample server celler server.
// @contact.name   Idris
// @contact.url    https://github.com/supuwoerc
// @contact.email  zhangzhouou@gmail.com
// @BasePath       /api/v1
func main() {
	if !isCli {
		defer bootstrap.Clean()
		bootstrap.Start()
	} else {
		cmd.Execute()
	}
}
