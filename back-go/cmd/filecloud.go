package main

import (
	"flag"
	"github.com/yddeng/filecloud/back-go/service"
)

var (
	file = flag.String("file", "./config.toml", "config file")
)

func main() {
	flag.Parse()

	logger := service.InitLogger("log", "filecloud")
	service.LoadConfig(*file)

	service.Launch()

	logger.Info("stop")
	service.Stop()

}
