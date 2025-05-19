package main

import (
	"github.com/supuwoerc/weaver/bootstrap"
	"github.com/supuwoerc/weaver/cmd"
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
