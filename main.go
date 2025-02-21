package main

import (
	"gin-web/bootstrap"
	"gin-web/cmd"
)

var isCli = false

func main() {
	defer bootstrap.Clean()
	if isCli {
		cmd.Execute()
	} else {
		bootstrap.Start()
	}
}
