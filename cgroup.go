package cgroup

import (
	"container/list"
	"runtime"
	"sync"
	"sync/atomic"
)

const (
	stopSubmitNo = 1 << iota
	stopSubmitYes
)

type Task func()

// CGroup represents a group of goroutines that can be executed concurrently.
type CGroup struct {
	concurrency int
	taskQueue   chan Task
	stop        int32
	wg          sync.WaitGroup
}

// New creates a new instance of CGroup with the given size.
// When concurrency is less than or equal to 0, it will be set to runtime.NumCPU() by default.
func New(concurrency int) *CGroup {
	if concurrency <= 0 {
		concurrency = runtime.NumCPU()
	}
	cg := &CGroup{
		concurrency: concurrency,
		taskQueue:   make(chan Task),
		wg:          sync.WaitGroup{},
		stop:        stopSubmitNo,
	}

	go cg.run()

	return cg
}

// run the goroutines that are waiting in the queue.
func (cg *CGroup) run() {
	taskList := list.New()

	stopSelect := false
	taskLimit := make(chan struct{}, cg.concurrency)
	cg.wg.Add(1)

	for {
		if !stopSelect {
			select {
			case task, ok := <-cg.taskQueue:
				if !ok {
					stopSelect = true
					continue
				}
				cg.wg.Add(1)
				taskList.PushBack(task)
			}
		}

		if cg.stop == stopSubmitYes && taskList.Len() == 0 && stopSelect {
			cg.wg.Done()
			return
		}

		for e := taskList.Front(); e != nil; e = e.Next() {
			taskLimit <- struct{}{}
			taskList.Remove(e)
			task := e.Value.(Task)
			go func() {
				defer func() {
					recover()
					<-taskLimit
					cg.wg.Done()
				}()
				task()
			}()
		}
	}
}

func (cg *CGroup) reset() {
	cg.taskQueue = make(chan Task)
	cg.wg = sync.WaitGroup{}
	go cg.run()
}

// Submit submits a task that needs to be executed.
func (cg *CGroup) Submit(task Task) {
	if cg.stop == stopSubmitYes {
		if atomic.CompareAndSwapInt32(&cg.stop, stopSubmitYes, stopSubmitNo) {
			cg.reset()
		}
	}
	cg.taskQueue <- task
}

// Wait blocks until all added tasks are executed.
func (cg *CGroup) Wait() {
	if atomic.CompareAndSwapInt32(&cg.stop, stopSubmitNo, stopSubmitYes) {
		close(cg.taskQueue)
		cg.wg.Wait()
		return
	}
}

// Async will cause previously added unfinished tasks to execute asynchronously without blocking waiting.
func (cg *CGroup) Async() {
	if atomic.CompareAndSwapInt32(&cg.stop, stopSubmitNo, stopSubmitYes) {
		close(cg.taskQueue)
		return
	}
}
