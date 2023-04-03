# cgroup
A more friendly implementation of limiting the number of concurrent coroutines

## Example:
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

	c := cgroup.New(size, func(i interface{}) {
		a := i.(int64)
		atomic.AddInt64(&sum, a)
	})
	for i := 0; i <= taskCount; i++ {
		c.Push(int64(i))
	}

	c.Wait()
	if sum != 5050 {
		panic("the value should equal 5050")
	}
}
```

## Installation
```
go get github.com/stoneflying/cgroup
```
