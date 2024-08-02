package main

import "gin-web/bootstrap"

func main() {
	defer bootstrap.Clean()
	bootstrap.Start()
}
