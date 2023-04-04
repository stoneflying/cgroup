# cgroup
simple and easy to concurrency limit in go

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
	
	for i := 1; i <= taskCount; i++ {
		a := int64(i)
		c.Submit(func() {
			atomic.AddInt64(&sum, a)
		})
	}

	c.Wait()
	if atomic.LoadInt64(&sum) != 5050 {
		t.Fatalf("the value should equal 5050, but got %v", sum)
	}
}
```

## installation
```
go get github.com/stoneflying/cgroup
```
