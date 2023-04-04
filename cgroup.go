package cgroup

import (
	"sync"
	"sync/atomic"
)

const (
	stopSubmitNo = 1 << iota
	stopSubmitYes
)

type Handle func()

// A CGroup instance represents a group of goroutine that
// can be executed concurrently
type CGroup struct {
	size int
	push chan Handle
	stop int32
	wg   sync.WaitGroup
}

// New create a concurrency limit instance.
// When size is less than or equal to zero, it means that there is no limit to the number of concurrency
func New(size int) *CGroup {
	c := &CGroup{
		size: size,
		push: make(chan Handle),
		wg:   sync.WaitGroup{},
		stop: stopSubmitNo,
	}

	go c.run()

	return c
}

func (c *CGroup) run() {
	buf := make([]Handle, 0, c.size)
	stopSelect := false
	ch := make(chan struct{}, c.size)

	for {
		if !stopSelect {
			select {
			case val, ok := <-c.push:
				if !ok {
					stopSelect = true
					continue
				}
				c.wg.Add(1)
				buf = append(buf, val)
			}
		}

		if c.stop == stopSubmitYes && len(buf) == 0 {
			return
		}

		for len(buf) > 0 {
			ch <- struct{}{}
			handle := buf[0]
			buf = buf[1:]

			go func() {
				defer func() {
					<-ch
					c.wg.Done()
				}()
				handle()
			}()
		}
	}
}

// Submit a fn that needs to be executed.
// When pushed a nil fn, after pushed data will be ignored.
func (c *CGroup) Submit(fn func()) {
	c.push <- fn
}

// Wait Block until all functions are executed.
func (c *CGroup) Wait() {
	if atomic.CompareAndSwapInt32(&c.stop, stopSubmitNo, stopSubmitYes) {
		close(c.push)
	}
	c.wg.Wait()
}
