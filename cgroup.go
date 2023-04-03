package cgroup

import (
	"sync"
	"sync/atomic"
)

const (
	pushEndNo = 1 << iota
	pushEndYes
)

// A CGroup instance represents a group of coroutines that
// can be executed concurrently
type CGroup struct {
	size    int
	fn      func(interface{})
	push    chan interface{}
	buf     []*interface{}
	wg      sync.WaitGroup
	pushEnd int32
	count   int
}

// New create a concurrency control instance to execute the specified fn.
// When size is less than or equal to zero, it means that there is no limit to the number of concurrency
func New(size int, fn func(interface{})) *CGroup {
	c := &CGroup{
		size:    size,
		fn:      fn,
		push:    make(chan interface{}),
		wg:      sync.WaitGroup{},
		buf:     make([]*interface{}, 0, size),
		pushEnd: pushEndNo,
		count:   0,
	}

	go c.run()

	return c
}

func (c *CGroup) run() {
	defer c.reset()

	for {
		if c.pushEnd != pushEndYes {
			select {
			case val, ok := <-c.push:
				if !ok {
					continue
				}
				c.wg.Add(1)
				c.buf = append(c.buf, &val)
			}
		}

		for (c.count < c.size || c.size <= 0) && len(c.buf) > 0 {
			c.count += 1
			v := c.buf[0]
			c.buf = c.buf[1:]
			go func() {
				defer func() {
					c.wg.Done()
					c.count -= 1
				}()

				c.fn(*v)
			}()
		}
	}
}

func (c *CGroup) reset() {
	c.push = nil
	c.size = 0
	c.fn = nil
	c.buf = nil
}

// Push a data that needs to be executed.
// Before you call the wait fn, you can always add data that needs to be executed.
// When the wait or pushEnd fn is called, the pushed data will be ignored.
func (c *CGroup) Push(data interface{}) {
	if c.pushEnd == pushEndYes {
		return
	}
	c.push <- data
}

// PushEnd indicates that there is no need to continue to add data that needs to be executed in the future.
// data added after calling this function will be ignored.
func (c *CGroup) PushEnd() {
	if atomic.CompareAndSwapInt32(&c.pushEnd, pushEndNo, pushEndYes) {
		close(c.push)
	}
}

// Wait blocks until all data added is executed by the specified fn
func (c *CGroup) Wait() {
	c.PushEnd()
	c.wg.Wait()
}
