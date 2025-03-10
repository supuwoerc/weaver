package main

import (
	"gin-web/bootstrap"
	"gin-web/cmd"
)

var isCli = false

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
