package cgroup

import (
	"container/list"
	"runtime"
	"sync"
	"sync/atomic"
)

type Task func()
type Logger interface {
	Printf(format string, args ...interface{})
}

const (
	OPENED = iota << 1
	CLOSED
)

// CGroup represents a group of goroutines that can be executed concurrently.
type CGroup struct {
	concurrency int
	push        chan Task
	wg          sync.WaitGroup
	options     *Options
	status      int32
	taskQueue   *list.List
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
		push:        make(chan Task),
		wg:          sync.WaitGroup{},
		options:     opts,
		status:      OPENED,
		taskQueue:   list.New(),
	}

	go cg.run()

	return cg
}

// run the goroutines that are waiting in the queue.
func (cg *CGroup) run() {
	stopSelect := false
	limit := make(chan struct{}, cg.concurrency)

	cg.wg.Add(1)

	for {
		if !stopSelect {
			select {
			case task, ok := <-cg.push:
				if !ok {
					stopSelect = true
					continue
				}
				cg.wg.Add(1)
				cg.taskQueue.PushBack(task)
			}
		}

		if cg.taskQueue.Len() == 0 && stopSelect {
			cg.wg.Done()
			return
		}

		for e := cg.taskQueue.Front(); e != nil; e = e.Next() {
			limit <- struct{}{}
			cg.taskQueue.Remove(e)
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
	close(cg.push)
}

// Submit submits a task that needs to be executed.
func (cg *CGroup) Submit(task Task) {
	if cg.IsClosed() {
		return
	}
	cg.push <- task
}

// Wait blocks until all added tasks are executed.
func (cg *CGroup) Wait() {
	if cg.IsClosed() {
		return
	}
	cg.Release()
	cg.wg.Wait()
}

// Release will cause previously added unfinished tasks to execute asynchronously without blocking waiting.
func (cg *CGroup) Release() {
	if !atomic.CompareAndSwapInt32(&cg.status, OPENED, CLOSED) {
		return
	}
	cg.reset()
}

func (cg *CGroup) IsClosed() bool {
	return atomic.LoadInt32(&cg.status) == CLOSED
}
