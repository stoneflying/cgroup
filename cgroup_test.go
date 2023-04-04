package cgroup

import (
	"sync/atomic"
	"testing"
)

func TestCGroupResult(t *testing.T) {
	sum := int64(0)
	size := 100
	taskCount := 100

	c := New(size)
	defer c.Release()

	for i := 0; i <= taskCount; i++ {
		a := int64(i)
		c.Submit(func() {
			atomic.AddInt64(&sum, a)
		})
	}

	c.Wait()
	if sum != 5050 {
		t.Fatalf("the value should equal 5050, but got %v", sum)
	}
}
