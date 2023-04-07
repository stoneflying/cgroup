package cgroup

import (
	"container/list"
	"runtime"
	"sync"
)

type Task func()
type Logger interface {
	Printf(format string, args ...interface{})
}

// CGroup represents a group of goroutines that can be executed concurrently.
type CGroup struct {
	concurrency int
	push        chan Task
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
		push:        make(chan Task),
		wg:          sync.WaitGroup{},
		options:     opts,
	}

	go cg.run()

	return cg
}

// run the goroutines that are waiting in the queue.
func (cg *CGroup) run() {
	taskQueue := list.New()

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
				taskQueue.PushBack(task)
			}
		}

		if taskQueue.Len() == 0 && stopSelect {
			cg.wg.Done()
			return
		}

		for e := taskQueue.Front(); e != nil; e = e.Next() {
			limit <- struct{}{}
			taskQueue.Remove(e)
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

// Submit submits a task that needs to be executed.
func (cg *CGroup) Submit(task Task) {
	cg.push <- task
}

// Wait blocks until all added tasks are executed.
func (cg *CGroup) Wait() {
	close(cg.push)
	cg.wg.Wait()
}

// Async will cause previously added unfinished tasks to execute asynchronously without blocking waiting.
func (cg *CGroup) Async() {
	close(cg.push)
}
