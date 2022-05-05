package task

import (
	"github.com/yddeng/utils"
	"runtime"
	"sync"
)

type Task interface {
	Do()
}

type funcTask struct {
	fn   interface{}
	args []interface{}
}

func (this *funcTask) Do() {
	_, _ = utils.CallFunc(this.fn, this.args...)
}

var (
	defaultTaskPool *TaskPool
	createOnce      sync.Once
)

func Default() *TaskPool {
	createOnce.Do(func() {
		defaultTaskPool = NewTaskPool(runtime.NumCPU()*2, defaultTaskSize)
	})
	return defaultTaskPool
}

func Submit(fn interface{}, args ...interface{}) error {
	return Default().Submit(fn, args...)
}

func SubmitTask(task Task, fullRet ...bool) error {
	return Default().SubmitTask(task, fullRet...)
}
