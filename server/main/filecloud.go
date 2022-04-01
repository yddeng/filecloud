package main

import (
	"filecloud/server"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("args: config")
	}

	server.InitLogger("log", "filecloud")
	server.LoadConfig(os.Args[1])

	server.Launch()
}
