package log

import (
	"fmt"
	"os"
)

/*
   全局默认提供一个Logger对外句柄，可以直接使用API系列调用
   不设置输出对象时，仅在控制台打印
   全局日志对象 stdOut
*/
var defaultLogger = newLogger(nil)

func Default() *Logger { return defaultLogger }

func CloseDebug() { defaultLogger.CloseDebug() }

func CloseStdOut() { defaultLogger.CloseStdOut() }

func SetOutput(basePath, fileName string, maxSize ...int) {
	defaultLogger.SetOutput(basePath, fileName, maxSize...)
}

func SetFlags(flag int) {
	defaultLogger.SetFlags(flag)
}

func SetPrefix(prefix string) {
	defaultLogger.SetPrefix(prefix)
}

func Debug(v ...interface{}) {
	if !defaultLogger.debugClosed {
		defaultLogger.output(DEBUG, "", v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if !defaultLogger.debugClosed {
		defaultLogger.output(DEBUG, format, v...)
	}
}

func Info(v ...interface{}) {
	defaultLogger.output(INFO, "", v...)
}

func Infof(format string, v ...interface{}) {
	defaultLogger.output(INFO, format, v...)
}

func Error(v ...interface{}) {
	defaultLogger.output(ERROR, "", v...)
}

func Errorf(format string, v ...interface{}) {
	defaultLogger.output(ERROR, format, v...)
}

func Fatal(v ...interface{}) {
	defaultLogger.output(FATAL, "", v...)
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	defaultLogger.output(FATAL, format, v...)
	os.Exit(1)
}

func Panic(v ...interface{}) {
	defaultLogger.output(PANIC, "", v...)
	panic(fmt.Sprintln(v...))
}

func Panicf(format string, v ...interface{}) {
	defaultLogger.output(PANIC, format, v...)
	panic(fmt.Sprintf(format, v...))
}

func Stack(v ...interface{}) {
	defaultLogger.output(ERROR, "", runStack(v...))
}
