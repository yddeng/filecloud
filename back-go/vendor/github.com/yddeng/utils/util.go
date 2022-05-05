package utils

import (
	"fmt"
	"runtime"
)

func Recover() (err error) {
	if r := recover(); r != nil {
		buf := make([]byte, 65535)
		l := runtime.Stack(buf, false)
		err = fmt.Errorf(fmt.Sprintf("%v: %s", r, buf[:l]))
	}
	return
}
