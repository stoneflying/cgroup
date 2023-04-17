package cgroup

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestCGroupResult(t *testing.T) {
	sum := int64(0)
	size := 50
	taskCount := 1000

	c := New(size)
	defer c.Release()

	for i := 1; i <= taskCount; i++ {
		a := int64(i)
		c.Submit(func() {
			atomic.AddInt64(&sum, a)
		})
	}

	c.Wait()
	if sum != 500500 {
		t.Fatalf("the value should equal 5050, but got %v", sum)
	}

	sum = 0
	c2 := New(size)
	defer c2.Release()
	for i := 1; i <= taskCount; i++ {
		a := int64(i)
		c2.Submit(func() {
			atomic.AddInt64(&sum, a)
		})
	}

	time.Sleep(3 * time.Second)
	if sum != 500500 {
		t.Fatalf("the value should equal 5050, but got %v", sum)
	}
}
