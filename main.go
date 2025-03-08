package main

import (
	"gin-web/bootstrap"
	"gin-web/cmd"
)

var isCli = false

func main() {
	if isCli {
		cmd.Execute()
	} else {
		defer bootstrap.Clean()
		bootstrap.Start()
	}
}
