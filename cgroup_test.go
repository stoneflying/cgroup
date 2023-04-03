package cgroup

import (
	"fmt"
	"sync/atomic"
	"testing"
)

func TestCGroupResult(t *testing.T) {
	sum := int64(0)
	size := 10
	taskCount := 100

	c := New(size)
	for i := 0; i <= taskCount; i++ {
		a := int64(i)
		c.Push(func() {
			atomic.AddInt64(&sum, a)
		})
	}

	c.Async()
	fmt.Println(sum)
	if sum != 5050 {
		t.Fatalf("the value should equal 5050, but got %v", sum)
	}
}
