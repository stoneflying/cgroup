package cgroup

import (
	"container/list"
	"runtime"
	"sync"
	"sync/atomic"
)

const (
	submitTaskFlag = 1 << iota
	syncWaitFlag
	asyncRunningFlag
)

type Task func()
type Logger interface {
	Printf(format string, args ...interface{})
}

// CGroup represents a group of goroutines that can be executed concurrently.
type CGroup struct {
	concurrency int
	taskQueue   chan Task
	status      int32
	wg          sync.WaitGroup
	options     *Options
}

// New creates a new instance of CGroup with the given size.
// When concurrency is less than or equal to 0, it will be set to runtime.NumCPU() by default.
func New(concurrency int, options ...Option) *CGroup {
	opts := loadOptions(options...)

	if concurrency <= 0 {
		concurrency = runtime.NumCPU()
	}

	cg := &CGroup{
		concurrency: concurrency,
		taskQueue:   make(chan Task),
		wg:          sync.WaitGroup{},
		status:      submitTaskFlag,
		options:     opts,
	}

	go cg.run()

	return cg
}

// run the goroutines that are waiting in the queue.
func (cg *CGroup) run() {
	taskList := list.New()

	stopSelect := false
	limit := make(chan struct{}, cg.concurrency)

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

		if (cg.status == syncWaitFlag || cg.status == asyncRunningFlag) && taskList.Len() == 0 && stopSelect {
			cg.wg.Done()
			return
		}

		for e := taskList.Front(); e != nil; e = e.Next() {
			limit <- struct{}{}
			taskList.Remove(e)
			task := e.Value.(Task)
			go func() {
				defer func() {
					if p := recover(); p != nil {
						if ph := cg.options.PanicHandler; ph != nil {
							ph(p)
						} else {
							cg.options.Logger.Printf("task panic stack begin: %v\n", p)
							var buf [4096]byte
							n := runtime.Stack(buf[:], false)
							cg.options.Logger.Printf("task panic stack end: %s\n", string(buf[:n]))
						}
					}
					<-limit
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
	if cg.status != submitTaskFlag {
		if atomic.CompareAndSwapInt32(&cg.status, cg.status&syncWaitFlag|cg.status&asyncRunningFlag, submitTaskFlag) {
			cg.reset()
		}
	}
	cg.taskQueue <- task
	return
}

// Wait blocks until all added tasks are executed.
func (cg *CGroup) Wait() {
	if atomic.CompareAndSwapInt32(&cg.status, submitTaskFlag, syncWaitFlag) {
		close(cg.taskQueue)
		cg.wg.Wait()
		return
	}
}

// Async will cause previously added unfinished tasks to execute asynchronously without blocking waiting.
func (cg *CGroup) Async() {
	if atomic.CompareAndSwapInt32(&cg.status, submitTaskFlag, asyncRunningFlag) {
		close(cg.taskQueue)
		return
	}
}
