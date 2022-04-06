package service

import "github.com/yddeng/utils/log"

var logger *log.Logger

func InitLogger(basePath string, fileName string) *log.Logger {
	logger = log.NewLogger(basePath, fileName, 1024*1024*4)
	logger.Debugf("%s logger init", fileName)
	return logger
}
