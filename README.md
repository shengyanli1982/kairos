English | [中文](./README_CN.md)

<div align="center" style="position:relative;">
	<img src="assets/logo.png" alt="logo">
</div>

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

		// 添加一个新的任务到调度器，并获取任务的 ID。
		// Add a new task to the scheduler and get the ID of the task.
		taskID := scheduler.Set("test", func(done ks.WaitForContextDone) (result any, err error) {
			// 当任务完成时，返回一个字符串。
			// When the task is done, return a string.
			for range done {
				return fmt.Sprintf("test task: %d", index), nil
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
# [CALLBACK] Task added, id: 1c4d337f-b266-4c5d-8f25-fb244ac82a66, name: test, execAt: 2024-03-29 15:22:41.71414 +0800 CST m=+0.200229378
% [MAIN] Task 0 can be retrieved, id: 1c4d337f-b266-4c5d-8f25-fb244ac82a66, name: test
# [CALLBACK] Task added, id: 1b884c0e-e5aa-41ba-8859-358673cbf602, name: test, execAt: 2024-03-29 15:22:41.714377 +0800 CST m=+0.200466654
% [MAIN] Task 1 can be retrieved, id: 1b884c0e-e5aa-41ba-8859-358673cbf602, name: test
# [CALLBACK] Task added, id: 4891cda4-35b0-41d7-a2c9-37b32caa42ef, name: test, execAt: 2024-03-29 15:22:41.714398 +0800 CST m=+0.200487260
% [MAIN] Task 2 can be retrieved, id: 4891cda4-35b0-41d7-a2c9-37b32caa42ef, name: test
# [CALLBACK] Task added, id: 807333f8-ec0d-42f7-88b9-20e08141da32, name: test, execAt: 2024-03-29 15:22:41.71441 +0800 CST m=+0.200499091
% [MAIN] Task 3 can be retrieved, id: 807333f8-ec0d-42f7-88b9-20e08141da32, name: test
# [CALLBACK] Task added, id: b34896fb-ec88-42f7-8141-0efa59184ecb, name: test, execAt: 2024-03-29 15:22:41.714428 +0800 CST m=+0.200517520
% [MAIN] Task 4 can be retrieved, id: b34896fb-ec88-42f7-8141-0efa59184ecb, name: test
# [CALLBACK] Task added, id: 0b6a4e10-7e52-4ee7-9236-d08a8c09cb7c, name: test, execAt: 2024-03-29 15:22:41.714443 +0800 CST m=+0.200532194
% [MAIN] Task 5 can be retrieved, id: 0b6a4e10-7e52-4ee7-9236-d08a8c09cb7c, name: test
# [CALLBACK] Task added, id: a1d4289f-5478-4573-bcdd-1a04cc779111, name: test, execAt: 2024-03-29 15:22:41.714454 +0800 CST m=+0.200543172
% [MAIN] Task 6 can be retrieved, id: a1d4289f-5478-4573-bcdd-1a04cc779111, name: test
# [CALLBACK] Task added, id: 62b2a14c-77f6-486a-b230-943c1b9c3818, name: test, execAt: 2024-03-29 15:22:41.714464 +0800 CST m=+0.200553433
% [MAIN] Task 7 can be retrieved, id: 62b2a14c-77f6-486a-b230-943c1b9c3818, name: test
# [CALLBACK] Task added, id: 613000e9-cb3f-4cfc-a21c-aee475b7e4f9, name: test, execAt: 2024-03-29 15:22:41.714487 +0800 CST m=+0.200576200
% [MAIN] Task 8 can be retrieved, id: 613000e9-cb3f-4cfc-a21c-aee475b7e4f9, name: test
# [CALLBACK] Task added, id: a873dd6a-0c5d-43f5-bd53-decce65e36cf, name: test, execAt: 2024-03-29 15:22:41.714531 +0800 CST m=+0.200620004
% [MAIN] Task 9 can be retrieved, id: a873dd6a-0c5d-43f5-bd53-decce65e36cf, name: test
# [CALLBACK] Task executed, id: 62b2a14c-77f6-486a-b230-943c1b9c3818, name: test, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task executed, id: 807333f8-ec0d-42f7-88b9-20e08141da32, name: test, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: 62b2a14c-77f6-486a-b230-943c1b9c3818, name: test
# [CALLBACK] Task removed, id: 807333f8-ec0d-42f7-88b9-20e08141da32, name: test
# [CALLBACK] Task executed, id: a1d4289f-5478-4573-bcdd-1a04cc779111, name: test, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task executed, id: 0b6a4e10-7e52-4ee7-9236-d08a8c09cb7c, name: test, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: a1d4289f-5478-4573-bcdd-1a04cc779111, name: test
# [CALLBACK] Task executed, id: a873dd6a-0c5d-43f5-bd53-decce65e36cf, name: test, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: a873dd6a-0c5d-43f5-bd53-decce65e36cf, name: test
# [CALLBACK] Task executed, id: 613000e9-cb3f-4cfc-a21c-aee475b7e4f9, name: test, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: 613000e9-cb3f-4cfc-a21c-aee475b7e4f9, name: test
# [CALLBACK] Task executed, id: 1c4d337f-b266-4c5d-8f25-fb244ac82a66, name: test, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: 1c4d337f-b266-4c5d-8f25-fb244ac82a66, name: test
# [CALLBACK] Task executed, id: 4891cda4-35b0-41d7-a2c9-37b32caa42ef, name: test, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: 4891cda4-35b0-41d7-a2c9-37b32caa42ef, name: test
# [CALLBACK] Task executed, id: b34896fb-ec88-42f7-8141-0efa59184ecb, name: test, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: b34896fb-ec88-42f7-8141-0efa59184ecb, name: test
# [CALLBACK] Task removed, id: 0b6a4e10-7e52-4ee7-9236-d08a8c09cb7c, name: test
# [CALLBACK] Task executed, id: 1b884c0e-e5aa-41ba-8859-358673cbf602, name: test, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: 1b884c0e-e5aa-41ba-8859-358673cbf602, name: test
```
