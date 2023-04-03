package cgroup

import (
	"sync"
)

type Handle func()

// A CGroup instance represents a group of coroutines that
// can be executed concurrently
type CGroup struct {
	size int
	push chan Handle
	wg   sync.WaitGroup
}

// New create a concurrency control instance to execute the specified fn.
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
				if val == nil {
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

// Push a data that needs to be executed.
// When pushed a nil fn, after pushed data will be ignored.
func (c *CGroup) Push(fn func()) *CGroup {
	c.push <- fn
	return c
}

// Async execute the fn, not will block the process.
func (c *CGroup) Async() {
	c.Push(nil)
}

// Wait blocks until all data added is executed by the specified fn
func (c *CGroup) Wait() {
	c.Push(nil)
	c.wg.Wait()
}
