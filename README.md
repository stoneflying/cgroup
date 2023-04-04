## cgroup
cgroup is a simple Go package that provides a way to execute a group of goroutines concurrently with a limited number of parallelism.

## Usage
#### 1.Import the package:
```
import "github.com/stoneflying/cgroup"
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

## Example
```
package main

import (
    "fmt"
    "github.com/username/cgroup"
    "time"
)

func main() {
    c := New(2)

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

## Output:
```
Task 1 started
Task 2 started
Task 2 finished
Task 1 finished
Task 3 started
Task 3 finished
All tasks finished
```