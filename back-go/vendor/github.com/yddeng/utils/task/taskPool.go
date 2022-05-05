package task

import (
	"errors"
	"sync"
)

const (
	defaultTaskSize = 1024
)

type TaskPool struct {
	workers    int
	workerSize int
	workerLock sync.Mutex

	taskChan chan Task

	die     chan struct{}
	dieOnce sync.Once
}

func (this *TaskPool) NumWorker() int {
	this.workerLock.Lock()
	defer this.workerLock.Unlock()
	return this.workers
}

func (this *TaskPool) NumTask() int {
	return len(this.taskChan)
}

func (this *TaskPool) submit(task Task, fullReturn bool) error {
	select {
	case <-this.die:
		return errors.New("taskPool:Submit pool is stopped")
	default:
	}

	var taskChan chan Task
	if this.workerSize == 0 {
		taskChan = make(chan Task, 1)
	} else {
		taskChan = this.taskChan
	}

	if fullReturn {
		select {
		case taskChan <- task:
		default:
			return errors.New("taskPool:Submit task channel is full")
		}
	} else {
		taskChan <- task
	}

	this.workerLock.Lock()
	defer this.workerLock.Unlock()

	if this.workerSize == 0 || this.workers < this.workerSize {
		this.workers++
		this.goWorker(taskChan)
	}
	return nil
}

func (this *TaskPool) goWorker(taskC chan Task) {
	go func() {
		defer func() {
			this.workerLock.Lock()
			this.workers--
			this.workerLock.Unlock()
		}()

		for {
			select {
			case task := <-taskC:
				task.Do()
			default:
				return
			}
		}
	}()

}

func (this *TaskPool) Submit(fn interface{}, args ...interface{}) error {
	return this.submit(&funcTask{fn: fn, args: args}, false)
}

func (this *TaskPool) SubmitTask(task Task, fullRet ...bool) error {
	fullReturn := false
	if len(fullRet) > 0 && fullRet[0] {
		fullReturn = true
	}
	return this.submit(task, fullReturn)
}

func (this *TaskPool) Stop() {
	this.dieOnce.Do(func() {
		close(this.die)
	})
}

// NewTaskPool
// workerSize > 0 , 限制goroutine的数量; workerSize = 0 , 不限制
func NewTaskPool(workerSize, taskSize int) *TaskPool {
	if taskSize < defaultTaskSize {
		taskSize = defaultTaskSize
	}
	if workerSize < 0 {
		workerSize = 0
	}

	pool := new(TaskPool)
	pool.die = make(chan struct{})
	pool.workerSize = workerSize
	pool.taskChan = make(chan Task, taskSize)

	return pool
}
