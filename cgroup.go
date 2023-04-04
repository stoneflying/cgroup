package cgroup

import (
	"sync"
)

type Handle func()

var releaseFn = Handle(func() {})

// A CGroup instance represents a group of goroutine that
// can be executed concurrently
type CGroup struct {
	size int
	push chan Handle
	wg   sync.WaitGroup
}

// New create a concurrency limit instance.
// When size is less than or equal to zero, it means that there is no limit to the number of concurrency
func New(size int) *CGroup {
	c := &CGroup{
		size: size,
		push: make(chan Handle),
		wg:   sync.WaitGroup{},
	}

	go c.run()

	return c
}

func (c *CGroup) run() {
	count := 0
	buf := make([]Handle, 0, c.size)
	stopPush := false

	for {
		if !stopPush {
			select {
			case val := <-c.push:
				if &val == &releaseFn {
					close(c.push)
					stopPush = true
					continue
				}
				c.wg.Add(1)
				buf = append(buf, val)
			}
		}

		for (count < c.size || c.size <= 0) && len(buf) > 0 {
			count += 1
			handle := buf[0]
			buf = buf[1:]
			go func() {
				defer func() {
					c.wg.Done()
					count -= 1
				}()
				handle()
			}()
		}
	}
}

// Submit a fn that needs to be executed.
func (c *CGroup) Submit(fn func()) *CGroup {
	c.push <- fn
	return c
}

// Release resources, you should always call this to avoid possible resource leaks.
func (c *CGroup) Release() {
	c.Submit(releaseFn)
}

// Wait Block until all functions are executed.
func (c *CGroup) Wait() {
	c.Submit(releaseFn)
	c.wg.Wait()
}
