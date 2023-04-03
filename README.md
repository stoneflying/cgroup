# cgroup
A more friendly implementation of limiting the number of concurrent coroutines

## basic example:
```
package main

import (
	"github.com/stoneflying/cgroup"
	"sync/atomic"
)

func main() {
	sum := int64(0)
	size := 10
	taskCount := 100

	c := cgroup.New(size)
	for i := 0; i <= taskCount; i++ {
		a := int64(i)
		c.Push(func() {
			atomic.AddInt64(&sum, a)
		})
	}

	c.Async()
	if sum != 5050 {
		panic("the value should equal 5050")
	}
}
```

## installation
```
go get github.com/stoneflying/cgroup
```
