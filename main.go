package main

import (
	"github.com/supuwoerc/weaver/bootstrap"
	"github.com/supuwoerc/weaver/cmd"
	_ "github.com/supuwoerc/weaver/docs"
)

var isCli = false

// @title			Weaver Service
// @version		1.0
// @description	Testing Swagger APIs.
// @termsOfService	https://github.com/supuwoerc/weaver
// @contact.name	API Support
// @contact.url	https://github.com/supuwoerc/weaver/issues
// @contact.email	zhangzhouou@gmain.com
// @license.name	MIT License
// @license.url	https://github.com/supuwoerc/weaver/blob/main/LICENSE
// @host			localhost:8804
// @BasePath		/api/v1
// @schemes		http https
func main() {
	switch {
	case isCli:
		cmd.Execute()
	default:
		app := bootstrap.WireApp()
		defer app.Close()
		app.Run()
	}
}
