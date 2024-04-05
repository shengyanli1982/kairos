English | [中文](./README_CN.md)

<div align="center" style="position:relative;">
	<img src="assets/logo.png" alt="logo">
</div>

[![Go Report Card](https://goreportcard.com/badge/github.com/shengyanli1982/kairos)](https://goreportcard.com/report/github.com/shengyanli1982/kairos)
[![Build Status](https://github.com/shengyanli1982/kairos/actions/workflows/test.yaml/badge.svg)](https://github.com/shengyanli1982/kairos/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/shengyanli1982/kairos.svg)](https://pkg.go.dev/github.com/shengyanli1982/kairos)

# Introduction

**Kairos** comes from the Greek word for "time". It means the right or opportune moment. It is a library equipped with a predefined timeout mechanism for executing specific tasks.

I wanted to simplify my development process and easily set up libraries to perform tasks at specific points in time. During my long-term development work, I realized that I needed to write a lot of code around context, which lacked generality. When something went wrong, it was difficult to pinpoint the source of errors. Much of the business code couldn't be reused across different projects.

# Why Kairos

To finish my work early and have more time to spend with my family. With this purpose in mind, I have abstracted tedious work tasks into a general function library for easy use. And of course, I want to help you too.

**The original intention of my design is:**

1. **Easy to use**: It doesn't take much time to learn.
2. **Highly reliable**: Use a small amount of code to complete complex tasks without introducing too many complex toolkits.
3. **Simple logic**: Use Golang's GMP coroutine model directly.

# Advantages

-   Simple and user-friendly
-   Lightweight with no external dependencies
-   Supports callback functions for custom actions

# Installation

```bash
go get github.com/shengyanli1982/kairos
```

# Quick Start

`Kairos` is very simple to use. Just few lines of code to get started.

## 1. Config

`Kairos` has a config object, which can be used to register callback functions. The config object has the following fields:

-   `WithCallback`: Register a callback function for the task.

```go
// Callback 是一个接口，定义了任务添加、执行和移除时的回调函数
// Callback is an interface that defines the callback functions when a task is added, executed, and removed
type Callback interface {
	// OnTaskAdded 是当任务被添加时的回调函数，它接收任务 id、任务名称和执行时间作为参数
	// OnTaskAdded is the callback function when a task is added, it takes the task id, task name, and execution time as parameters
	OnTaskAdded(id, name string, execAt time.Time)

	// OnTaskExecuted 是当任务被执行时的回调函数，它接收任务 id、任务名称、数据、原因和错误作为参数
	// OnTaskExecuted is the callback function when a task is executed, it takes the task id, task name, data, reason, and error as parameters
	OnTaskExecuted(id, name string, data interface{}, reason, err error)

	// OnTaskRemoved 是当任务被移除时的回调函数，它接收任务 id 和任务名称作为参数
	// OnTaskRemoved is the callback function when a task is removed, it takes the task id and task name as parameters
	OnTaskRemoved(id, name string)
}
```

## 2. Methods

The `Kairos` provides the following methods:

-   `New`: Create a new `Scheduler` object. The `Scheduler` object is used to manage tasks.
-   `Stop`: Stop the `Scheduler`. If the `Scheduler` object is stopped, all tasks will be stopped and removed.
-   `Set`: Add a task to the `Scheduler`. The `Set` method takes the task `name`, the `delay` time.Duration to execute the task, and `handleFunc` to the task as parameters.
-   `SetAt`: Add a task to the `Scheduler` at a specific time. The `SetAt` method takes the task `name`, the `execAt` time.Time to execute the task, and `handleFunc` to the task as parameters.
-   `Get`: Get the task from the `Scheduler` by the task `id`.
-   `Delete`: Delete the task from the `Scheduler` by the task `id`.

## 3. Task

The `Task` is a crucial concept in `Kairos`. It allows for the execution of specific tasks at designated times. The `Task` object provides the following methods:

-   `GetMetadata`: Retrieves the metadata of the task, which includes methods to obtain task information.
    1.  `GetID`: Retrieves the task ID.
    2.  `GetName`: Retrieves the task name.
    3.  `GetHandleFunc`: Retrieves the task handle function.
-   `EarlyReturn`: Manually stops task execution and returns early, without waiting for the timeout or cancel signal. It invokes the `handleFunc`.
-   `Cancel`: Manually stops task execution and returns immediately, without executing the `handleFunc`.
-   `Wait`: Waits for the task to complete, blocking the current goroutine until the task is finished.

> [!NOTE]
>
> The `Wait` method is a blocking method. It is recommended to use it in a separate `goroutine` to avoid blocking the main `goroutine`.
>
> If you need to retrieve a task, you can use the `Get` method to obtain the task by its `ID`. Once you have the task, you can use the `Wait` method to wait for its completion.
>
> There is no need to manually stop a task. After the task is executed, its status is automatically cleared in the `Scheduler`. If you want to actively delete a task, you can use the `Delete` method to remove the corresponding task (this should only be done if necessary).

## 4. Example

Example code is located in the `examples` directory.

### 4.1 Simple Example

With just a few lines of code, you can effortlessly create and run a task while obtaining its information.

```go
package main

import (
	"fmt"
	"time"

	ks "github.com/shengyanli1982/kairos"
)

func main() {
	// 创建一个新的配置，并设置回调函数。
	// Create a new configuration and set the callback function.
	config := ks.NewConfig()

	// 使用配置创建一个新的调度器。
	// Create a new scheduler with the configuration.
	scheduler := ks.New(config)

	// 停止调度器，停止所有的任务。
	// Stop the scheduler, stop all tasks.
	scheduler.Stop()

	// 添加一个新的任务到调度器，并获取任务的 ID。
	// Add a new task to the scheduler and get the ID of the task.
	taskID := scheduler.Set("test_task", func(done ks.WaitForContextDone) (result any, err error) {
		// 做你想做的任何事情
		// Do whatever you want to do
		// 这里我们模拟一个需要 100 毫秒才能完成的任务。
		// Here we simulate a task that takes 100 milliseconds to complete.
		time.Sleep(time.Millisecond * 100)

		// 如果任务没有完成，返回 nil。
		// If the task is not done, return nil.
		return nil, nil

		// 设置任务的延迟时间为 200 毫秒。
		// Set the delay time of the task to 200 milliseconds.
	}, time.Millisecond*200)

	// 获取任务。
	// Get the task.
	task := scheduler.Get(taskID)

	// 打印任务添加的信息。
	// Print the information when the task is added.
	fmt.Printf("%% [MAIN] Task can be retrieved, id: %s, name: %s\n", task.GetMetadata().GetID(), task.GetMetadata().GetName())
}
```

**Result**

```bash
$ go run demo.go
% [MAIN] Task can be retrieved, id: 045a6ecb-5d4d-4393-84a3-679c943bd6f7, name: test_task
```

### 4.2 Full Example

This example showcases the usage of a callback function to handle task execution, addition, and removal.

```go
package main

import (
	"fmt"
	"time"

	ks "github.com/shengyanli1982/kairos"
)

// demoSchedCallback 结构体实现了调度器的回调接口。
// The demoSchedCallback struct implements the callback interface of the scheduler.
type demoSchedCallback struct{}

// OnTaskExecuted 是一个方法，当任务执行后会被调用。
// OnTaskExecuted is a method that is called after a task is executed.
func (tc *demoSchedCallback) OnTaskExecuted(id, name string, data any, reason, err error) {
	// 打印任务执行后的信息。
	// Print the information after the task is executed.
	fmt.Printf("# [CALLBACK] Task executed, id: %s, name: %s, data: %v, reason: %v, err: %v\n", id, name, data, reason, err)
}

// OnTaskAdded 是一个方法，当任务被添加时会被调用。
// OnTaskAdded is a method that is called when a task is added.
func (tc *demoSchedCallback) OnTaskAdded(id, name string, execAt time.Time) {
	// 打印任务添加的信息。
	// Print the information when the task is added.
	fmt.Printf("# [CALLBACK] Task added, id: %s, name: %s, execAt: %s\n", id, name, execAt.String())
}

// OnTaskRemoved 是一个方法，当任务被删除时会被调用。
// OnTaskRemoved is a method that is called when a task is removed.
func (tc *demoSchedCallback) OnTaskRemoved(id, name string) {
	// 打印任务删除的信息。
	// Print the information when the task is removed.
	fmt.Printf("# [CALLBACK] Task removed, id: %s, name: %s\n", id, name)
}

// main 是程序的入口点。
// main is the entry point of the program.
func main() {
	// 创建一个新的配置，并设置回调函数。
	// Create a new configuration and set the callback function.
	config := ks.NewConfig().WithCallback(&demoSchedCallback{})

	// 使用配置创建一个新的调度器。
	// Create a new scheduler with the configuration.
	scheduler := ks.New(config)

	// 循环添加 10 个任务到调度器。
	// Loop to add 10 tasks to the scheduler.
	for i := 0; i < 10; i++ {
		// 保存当前的索引值。
		// Save the current index value.
		index := i

		// taskName 是任务的名称，通过格式化字符串生成
		// taskName is the name of the task, generated by formatting a string
		taskName := fmt.Sprintf("test_task_%d", index)

		// 添加一个新的任务到调度器，并获取任务的 ID。
		// Add a new task to the scheduler and get the ID of the task.
		taskID := scheduler.Set(taskName, func(done ks.WaitForContextDone) (result any, err error) {

			// 当任务完成时，返回任务的名称 (这步不是必须，只是为了表示函数内部可以接收外部的 ctx 信号)。
			// When the task is done, return the name of the task (this step is not necessary, just to indicate that the function can receive the ctx signal from the outside).
			for range done {
				return taskName, nil
			}

			// 如果任务没有完成，返回 nil。
			// If the task is not done, return nil.
			return nil, nil

			// 设置任务的延迟时间为 200 毫秒。
			// Set the delay time of the task to 200 milliseconds.
		}, time.Millisecond*200)

		// 获取任务。
		// Get the task.
		task := scheduler.Get(taskID)

		// 打印任务添加的信息。
		// Print the information when the task is added.
		fmt.Printf("%% [MAIN] Task %d can be retrieved, id: %s, name: %s\n", index, task.GetMetadata().GetID(), task.GetMetadata().GetName())
	}

	// 等待一段时间，让任务有机会执行。
	// Wait for a while to give the tasks a chance to execute.
	time.Sleep(time.Millisecond * 500)

	// 停止调度器，停止所有的任务。
	// Stop the scheduler, stop all tasks.
	scheduler.Stop()
}
```

**Result**

```bash
$ go run demo.go
# [CALLBACK] Task added, id: 63f053cf-5da8-4070-802c-2d9c9f493c23, name: test_task_0, execAt: 2024-03-29 22:14:20.347086 +0800 CST m=+0.200244340
% [MAIN] Task 0 can be retrieved, id: 63f053cf-5da8-4070-802c-2d9c9f493c23, name: test_task_0
# [CALLBACK] Task added, id: 3006c080-ac62-4523-9ff7-a7d21b5d6e23, name: test_task_1, execAt: 2024-03-29 22:14:20.347357 +0800 CST m=+0.200515252
% [MAIN] Task 1 can be retrieved, id: 3006c080-ac62-4523-9ff7-a7d21b5d6e23, name: test_task_1
# [CALLBACK] Task added, id: 585408bd-1f31-433e-be4b-a7f03e4ccb58, name: test_task_2, execAt: 2024-03-29 22:14:20.347379 +0800 CST m=+0.200536895
% [MAIN] Task 2 can be retrieved, id: 585408bd-1f31-433e-be4b-a7f03e4ccb58, name: test_task_2
# [CALLBACK] Task added, id: 6ea38719-503e-4fce-9155-03280a818c0d, name: test_task_3, execAt: 2024-03-29 22:14:20.347391 +0800 CST m=+0.200548818
% [MAIN] Task 3 can be retrieved, id: 6ea38719-503e-4fce-9155-03280a818c0d, name: test_task_3
# [CALLBACK] Task added, id: 8d7a4d5e-27b9-47c7-b307-445e84e8fa79, name: test_task_4, execAt: 2024-03-29 22:14:20.34741 +0800 CST m=+0.200567656
% [MAIN] Task 4 can be retrieved, id: 8d7a4d5e-27b9-47c7-b307-445e84e8fa79, name: test_task_4
# [CALLBACK] Task added, id: 0ebb9759-18ea-4476-8b78-d1840a547ef2, name: test_task_5, execAt: 2024-03-29 22:14:20.347428 +0800 CST m=+0.200586399
% [MAIN] Task 5 can be retrieved, id: 0ebb9759-18ea-4476-8b78-d1840a547ef2, name: test_task_5
# [CALLBACK] Task added, id: 8234e5a4-e06d-4bc6-8642-949ae86ac2e6, name: test_task_6, execAt: 2024-03-29 22:14:20.347439 +0800 CST m=+0.200597544
% [MAIN] Task 6 can be retrieved, id: 8234e5a4-e06d-4bc6-8642-949ae86ac2e6, name: test_task_6
# [CALLBACK] Task added, id: 8c2f8b22-3fd6-4aff-97b2-07c6b18aeb43, name: test_task_7, execAt: 2024-03-29 22:14:20.34745 +0800 CST m=+0.200607930
% [MAIN] Task 7 can be retrieved, id: 8c2f8b22-3fd6-4aff-97b2-07c6b18aeb43, name: test_task_7
# [CALLBACK] Task added, id: c3758e02-7cd5-45ce-8618-d83a3843c5a1, name: test_task_8, execAt: 2024-03-29 22:14:20.34746 +0800 CST m=+0.200618501
% [MAIN] Task 8 can be retrieved, id: c3758e02-7cd5-45ce-8618-d83a3843c5a1, name: test_task_8
# [CALLBACK] Task added, id: a7f95604-fe2b-46b1-a762-44b69228ed43, name: test_task_9, execAt: 2024-03-29 22:14:20.34749 +0800 CST m=+0.200648064
% [MAIN] Task 9 can be retrieved, id: a7f95604-fe2b-46b1-a762-44b69228ed43, name: test_task_9
# [CALLBACK] Task executed, id: 8d7a4d5e-27b9-47c7-b307-445e84e8fa79, name: test_task_4, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task executed, id: 3006c080-ac62-4523-9ff7-a7d21b5d6e23, name: test_task_1, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: 8d7a4d5e-27b9-47c7-b307-445e84e8fa79, name: test_task_4
# [CALLBACK] Task removed, id: 3006c080-ac62-4523-9ff7-a7d21b5d6e23, name: test_task_1
# [CALLBACK] Task executed, id: 63f053cf-5da8-4070-802c-2d9c9f493c23, name: test_task_0, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task executed, id: 0ebb9759-18ea-4476-8b78-d1840a547ef2, name: test_task_5, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: 0ebb9759-18ea-4476-8b78-d1840a547ef2, name: test_task_5
# [CALLBACK] Task executed, id: 8c2f8b22-3fd6-4aff-97b2-07c6b18aeb43, name: test_task_7, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: 8c2f8b22-3fd6-4aff-97b2-07c6b18aeb43, name: test_task_7
# [CALLBACK] Task executed, id: a7f95604-fe2b-46b1-a762-44b69228ed43, name: test_task_9, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: a7f95604-fe2b-46b1-a762-44b69228ed43, name: test_task_9
# [CALLBACK] Task executed, id: 585408bd-1f31-433e-be4b-a7f03e4ccb58, name: test_task_2, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task executed, id: c3758e02-7cd5-45ce-8618-d83a3843c5a1, name: test_task_8, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task executed, id: 8234e5a4-e06d-4bc6-8642-949ae86ac2e6, name: test_task_6, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: 8234e5a4-e06d-4bc6-8642-949ae86ac2e6, name: test_task_6
# [CALLBACK] Task removed, id: 585408bd-1f31-433e-be4b-a7f03e4ccb58, name: test_task_2
# [CALLBACK] Task removed, id: c3758e02-7cd5-45ce-8618-d83a3843c5a1, name: test_task_8
# [CALLBACK] Task removed, id: 63f053cf-5da8-4070-802c-2d9c9f493c23, name: test_task_0
# [CALLBACK] Task executed, id: 6ea38719-503e-4fce-9155-03280a818c0d, name: test_task_3, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: 6ea38719-503e-4fce-9155-03280a818c0d, name: test_task_3
```
