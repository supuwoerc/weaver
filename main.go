package main

import "gin-web/cmd"

func main() {
	defer cmd.Clean()
	cmd.Start()
}
