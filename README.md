## cgroup
CGroup is a simple Go library that provides an easy way to execute tasks concurrently with controlled concurrency. It is designed to help manage and limit the number of concurrent tasks that are executed in a Go program.

## Features
#### 1.Ability to control the maximum number of tasks running concurrently
#### 2.Asynchronous execution of tasks
#### 3.Simple and easy to use API
#### 4.non-blocking cubmit task

## Getting Started
#### 1.Installation:
```
go get github.com/stoneflying/cgroup
```

#### 2.Create a new CGroup with the desired number of parallelism:
```
// Create a new CGroup with 5 parallelism
c := cgroup.New(5)
```

#### 3.Submit the tasks that need to be executed:
```
// Submit a task
c.Submit(func() {
    // code to be executed concurrently
})
```

#### 4.Wait for all tasks to be executed:
```
// Wait for all tasks to complete
c.Wait()
```

#### 5.Optionally, if you need to execute the tasks asynchronously without blocking, you can use the Async method:
```
// Execute the tasks asynchronously
c.Async()
```

## Basic example
```
package main

import (
	"fmt"
	"github.com/stoneflying/cgroup"
	"time"
)

func main() {
	c := cgroup.New(2)
	c.Submit(func() {
		fmt.Println("Task 1 started")
		time.Sleep(2 * time.Second)
		fmt.Println("Task 1 finished")
	})

	c.Submit(func() {
		fmt.Println("Task 2 started")
		time.Sleep(1 * time.Second)
		fmt.Println("Task 2 finished")
	})

	c.Submit(func() {
		fmt.Println("Task 3 started")
		time.Sleep(3 * time.Second)
		fmt.Println("Task 3 finished")
	})

	c.Wait()
	fmt.Println("All tasks finished")
}
```

### Output:
```
Task 1 started
Task 2 started
Task 2 finished
Task 3 started
Task 1 finished
Task 3 finished
All tasks finished
```

## Attention
#### 1.You should always call await or async func at last, otherwise there will be resource leaks.
#### 2.When you called async or await func, you should not continue to add execution func.
#### 3.Don't repeat call async or await func.

## Contributing
Contributions are welcome!   
For bug reports or feature requests, please open an issue.   
For pull requests, please make sure your changes are covered by tests and 
follow the same coding standards as the existing code.
