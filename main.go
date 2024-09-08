package main

import (
	"gin-web/bootstrap"
	"gin-web/cmd"
)

var isCli = false

func main() {
	if !isCli {
		defer bootstrap.Clean()
		bootstrap.Start()
	} else {
		cmd.Execute()
	}
}
