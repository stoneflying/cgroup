package cgroup

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestCGroupResult(t *testing.T) {
	sum := int64(0)
	size := 50
	taskCount := 100

	c := New(size)

	for i := 1; i <= taskCount; i++ {
		a := int64(i)
		c.Submit(func() {
			atomic.AddInt64(&sum, a)
		})
	}

	c.Wait()
	if sum != 5050 {
		t.Fatalf("the value should equal 5050, but got %v", sum)
	}

	for i := 1; i <= taskCount; i++ {
		a := int64(i)
		c.Submit(func() {
			atomic.AddInt64(&sum, a)
		})
	}
	c.Wait()
	if sum != 10100 {
		t.Fatalf("the value should equal 10100, but got %v", sum)
	}

	for i := 1; i <= taskCount; i++ {
		a := int64(i)
		c.Submit(func() {
			atomic.AddInt64(&sum, a)
		})
	}
	c.Async()

	time.Sleep(3 * time.Second)
	if sum != 15150 {
		t.Fatalf("the value should equal 10150, but got %v", sum)
	}
}
