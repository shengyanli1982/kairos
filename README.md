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

    	// OnTaskExecuted 是当任务被执行时的回调函数，它接收任务 id、任务名称、任务结果、原因和错误作为参数
    	// OnTaskExecuted is the callback function when a task is executed, it takes the task id, task name, task result, reason, and error as parameters
    	OnTaskExecuted(id, name string, result interface{}, reason, err error)

    	// OnTaskRemoved 是当任务被移除时的回调函数，它接收任务 id 和任务名称作为参数
    	// OnTaskRemoved is the callback function when a task is removed, it takes the task id and task name as parameters
    	OnTaskRemoved(id, name string)

    	// OnTaskDuplicated 是当任务重复时的回调函数，它接收任务 id 和任务名称作为参数
    	// OnTaskDuplicated is the callback function when a task is duplicated, it takes the task id and task name as parameters
    	OnTaskDuplicated(id, name string)
    }
    ```

-   `WithUniqued`: Disable duplicated tasks. When set to `true`, the `Scheduler` will not allow tasks with the same name.

## 2. Methods

The `Kairos` provides the following methods:

-   `New`: Create a new `Scheduler` object. The `Scheduler` object is used to manage tasks.
-   `Stop`: Stop the `Scheduler`. If the `Scheduler` object is stopped, all tasks will be stopped and removed.
-   `Set`: Add a task to the `Scheduler`. The `Set` method takes the task `name`, the `delay` time.Duration to execute the task, and `handleFunc` to the task as parameters.
-   `SetAt`: Add a task to the `Scheduler` at a specific time. The `SetAt` method takes the task `name`, the `execAt` time.Time to execute the task, and `handleFunc` to the task as parameters.
-   `Get`: Get the task from the `Scheduler` by the task `id`.
-   `Delete`: Delete the task from the `Scheduler` by the task `id`.
-   `Count`: Retrieve the number of tasks in the `Scheduler`.

> [!TIP]
>
> If you want to ensure that tasks in the `Scheduler` are unique, you can use the `WithUniqued` option to disable duplicated tasks.
>
> By default, the `Scheduler` allows duplicated tasks (`WithUniqued` is set to `false`). This means you can add multiple tasks with the same name to the `Scheduler`.
>
> If you want the `Scheduler` to use the task name as the identifier, make sure to use different names for each task when `WithUniqued` is set to `true`.
>
> When `WithUniqued` is set to `true`, the `Set` and `SetAt` methods will return the `id` of the running task.
>
> If you want to access the result value of a custom handler set by `Set` and `SetAt`, you can utilize the `OnTaskExecuted` method in `Callback`. This method has a `result` parameter, which represents the result value returned by the custom handler.

## 3. Task

The `Task` is a crucial concept in `Kairos`. It allows for the execution of specific tasks at designated times. The `Task` object provides the following methods:

-   `GetMetadata`: Retrieves the metadata of the task, which includes methods to obtain task information.
    1.  `GetID`: Retrieves the task `id`.
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
	defer scheduler.Stop()

	// 添加一个新的任务到调度器，并获取任务的 ID。
	// Add a new task to the scheduler and get the ID of the task.
	taskID, err := scheduler.Set("test_task", func(done ks.WaitForContextDone) (result any, err error) {
		// 这是任务的主体部分，你可以在这里做你想做的任何事情。
		// This is the main part of the task, you can do whatever you want to do here.
		// 这里我们只是让任务休眠 100 毫秒，模拟任务的执行。
		// Here we just let the task sleep for 100 milliseconds, simulating the execution of the task.
		time.Sleep(time.Millisecond * 100)

		// 如果任务没有完成，返回 nil。
		// If the task is not done, return nil.
		return nil, nil
		// 设置任务的延迟时间为 200 毫秒。
		// Set the delay time of the task to 200 milliseconds.
	}, time.Millisecond*200)

	// 如果设置任务时出现错误
	// If an error occurs when setting the task
	if err != nil {
		// 打印错误信息。
		// Print the error message.
		fmt.Printf("%% [MAIN] Error: %s\n", err.Error())

		// 返回主函数，不执行后续的代码。
		// Return to the main function, do not execute the following code.
		return
	}

	// 使用任务的 ID 获取任务。
	// Get the task using the ID of the task.
	task, err := scheduler.Get(taskID)

	// 如果获取任务时出现错误
	// If an error occurs when getting the task
	if err != nil {
		// 打印错误信息。
		// Print the error message.
		fmt.Printf("%% [MAIN] Error: %s\n", err.Error())

		// 返回主函数，不执行后续的代码。
		// Return to the main function, do not execute the following code.
		return
	}

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
func (tc *demoSchedCallback) OnTaskExecuted(id, name string, result any, reason, err error) {
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

// OnTaskDuplicated 是一个方法，当任务重复时会被调用。
// OnTaskDuplicated is a method that is called when a task is duplicated.
func (tc *demoSchedCallback) OnTaskDuplicated(id, name string) {
	// 打印任务重复的信息。
	// Print the information when the task is duplicated.
	fmt.Printf(">>> [CALLBACK] Task duplicated , id: %s, name: %s\n", id, name)
}

// main 是程序的入口点。
// main is the entry point of the program.
func main() {
	// 创建一个新的配置，并设置回调函数。
	// Create a new configuration and set the callback function.
	config := ks.NewConfig().WithCallback(&demoSchedCallback{}).WithUniqued(true)

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
		taskID, err := scheduler.Set(taskName, func(done ks.WaitForContextDone) (result any, err error) {

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

		// 如果设置任务时出现错误
		// If an error occurs when setting the task
		if err != nil {
			// 打印错误信息。
			// Print the error message.
			fmt.Printf("%% [MAIN] Error: %s\n", err.Error())

			// 返回主函数，不执行后续的代码。
			// Return to the main function, do not execute the following code.
			return
		}

		// 获取任务。
		// Get the task.
		task, err := scheduler.Get(taskID)

		// 如果获取任务时出现错误
		// If an error occurs when getting the task
		if err != nil {
			// 打印错误信息。
			// Print the error message.
			fmt.Printf("%% [MAIN] Error: %s\n", err.Error())

			// 返回主函数，不执行后续的代码。
			// Return to the main function, do not execute the following code.
			return
		}

		// 打印任务添加的信息。
		// Print the information when the task is added.
		fmt.Printf("%% [MAIN] Task %d can be retrieved, id: %s, name: %s\n", index, task.GetMetadata().GetID(), task.GetMetadata().GetName())
	}

	// 添加一个名字存在的任务到调度器，等待触发回调。
	// Add a task with an existing name to the scheduler and wait for the callback to be triggered.
	repeatTaskName := "test_task_9"
	taskID, err := scheduler.Set(repeatTaskName, func(done ks.WaitForContextDone) (result any, err error) {
		// 当任务完成时，返回任务的名称 (这步不是必须，只是为了表示函数内部可以接收外部的 ctx 信号)。
		// When the task is done, return the name of the task (this step is not necessary, just to indicate that the function can receive the ctx signal from the outside).
		for range done {
			return repeatTaskName, nil
		}

		// 如果任务没有完成，返回 nil。
		// If the task is not done, return nil.
		return nil, nil

		// 设置任务的延迟时间为 200 毫秒。
		// Set the delay time of the task to 200 milliseconds.
	}, time.Millisecond*200)

	// 如果设置任务时出现错误
	// If an error occurs when setting the task
	if err != nil {
		// 打印错误信息。
		// Print the error message.
		fmt.Printf("%% [MAIN] Error: %s\n", err.Error())

		// 返回主函数，不执行后续的代码。
		// Return to the main function, do not execute the following code.
		return
	}

	// 打印任务添加的信息。
	// Print the information when the task is added.
	fmt.Printf("%% [MAIN] The duplicate task can be retrieved, id: %s, name: %s\n", taskID, repeatTaskName)

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
# [CALLBACK] Task added, id: dddac11f-4c74-4708-80bf-989f510b6cad, name: test_task_0, execAt: 2024-04-06 16:49:53.798281 +0800 CST m=+0.200312335
% [MAIN] Task 0 can be retrieved, id: dddac11f-4c74-4708-80bf-989f510b6cad, name: test_task_0
# [CALLBACK] Task added, id: db88bbcc-9235-40ca-800c-a3ec42d1c2f2, name: test_task_1, execAt: 2024-04-06 16:49:53.798552 +0800 CST m=+0.200583830
% [MAIN] Task 1 can be retrieved, id: db88bbcc-9235-40ca-800c-a3ec42d1c2f2, name: test_task_1
# [CALLBACK] Task added, id: e518d082-5a2a-4028-b102-d41f87f8caf8, name: test_task_2, execAt: 2024-04-06 16:49:53.798575 +0800 CST m=+0.200606459
% [MAIN] Task 2 can be retrieved, id: e518d082-5a2a-4028-b102-d41f87f8caf8, name: test_task_2
# [CALLBACK] Task added, id: 74cab4d3-0e2b-4489-a8d5-4b5ba071daf5, name: test_task_3, execAt: 2024-04-06 16:49:53.798587 +0800 CST m=+0.200618709
% [MAIN] Task 3 can be retrieved, id: 74cab4d3-0e2b-4489-a8d5-4b5ba071daf5, name: test_task_3
# [CALLBACK] Task added, id: 2b98b896-822d-4b2a-91cd-3b74cfdd7acd, name: test_task_4, execAt: 2024-04-06 16:49:53.798607 +0800 CST m=+0.200638692
% [MAIN] Task 4 can be retrieved, id: 2b98b896-822d-4b2a-91cd-3b74cfdd7acd, name: test_task_4
# [CALLBACK] Task added, id: ebb40945-2bfd-494d-933a-ced6a2f85c3d, name: test_task_5, execAt: 2024-04-06 16:49:53.798649 +0800 CST m=+0.200681205
% [MAIN] Task 5 can be retrieved, id: ebb40945-2bfd-494d-933a-ced6a2f85c3d, name: test_task_5
# [CALLBACK] Task added, id: a498f460-392b-4ffe-a905-7689e337eeb8, name: test_task_6, execAt: 2024-04-06 16:49:53.798676 +0800 CST m=+0.200707531
% [MAIN] Task 6 can be retrieved, id: a498f460-392b-4ffe-a905-7689e337eeb8, name: test_task_6
# [CALLBACK] Task added, id: ea23211f-0961-49df-804e-83c93d80ce85, name: test_task_7, execAt: 2024-04-06 16:49:53.798708 +0800 CST m=+0.200739353
% [MAIN] Task 7 can be retrieved, id: ea23211f-0961-49df-804e-83c93d80ce85, name: test_task_7
# [CALLBACK] Task added, id: be65abf3-e1d5-403c-a180-05cb42a80fb0, name: test_task_8, execAt: 2024-04-06 16:49:53.798732 +0800 CST m=+0.200763759
% [MAIN] Task 8 can be retrieved, id: be65abf3-e1d5-403c-a180-05cb42a80fb0, name: test_task_8
# [CALLBACK] Task added, id: f0f3d67b-7563-4c1a-b004-5e96beca0be4, name: test_task_9, execAt: 2024-04-06 16:49:53.798868 +0800 CST m=+0.200899985
% [MAIN] Task 9 can be retrieved, id: f0f3d67b-7563-4c1a-b004-5e96beca0be4, name: test_task_9
>>> [CALLBACK] Task duplicated , id: f0f3d67b-7563-4c1a-b004-5e96beca0be4, name: test_task_9
# [CALLBACK] Task added, id: f0f3d67b-7563-4c1a-b004-5e96beca0be4, name: test_task_9, execAt: 2024-04-06 16:49:53.798893 +0800 CST m=+0.200925081
% [MAIN] The duplicate task can be retrieved, id: f0f3d67b-7563-4c1a-b004-5e96beca0be4, name: test_task_9
# [CALLBACK] Task executed, id: dddac11f-4c74-4708-80bf-989f510b6cad, name: test_task_0, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task executed, id: e518d082-5a2a-4028-b102-d41f87f8caf8, name: test_task_2, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: dddac11f-4c74-4708-80bf-989f510b6cad, name: test_task_0
# [CALLBACK] Task removed, id: e518d082-5a2a-4028-b102-d41f87f8caf8, name: test_task_2
# [CALLBACK] Task executed, id: 2b98b896-822d-4b2a-91cd-3b74cfdd7acd, name: test_task_4, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task executed, id: f0f3d67b-7563-4c1a-b004-5e96beca0be4, name: test_task_9, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: f0f3d67b-7563-4c1a-b004-5e96beca0be4, name: test_task_9
# [CALLBACK] Task executed, id: ebb40945-2bfd-494d-933a-ced6a2f85c3d, name: test_task_5, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: ebb40945-2bfd-494d-933a-ced6a2f85c3d, name: test_task_5
# [CALLBACK] Task executed, id: a498f460-392b-4ffe-a905-7689e337eeb8, name: test_task_6, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: a498f460-392b-4ffe-a905-7689e337eeb8, name: test_task_6
# [CALLBACK] Task executed, id: 74cab4d3-0e2b-4489-a8d5-4b5ba071daf5, name: test_task_3, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: 74cab4d3-0e2b-4489-a8d5-4b5ba071daf5, name: test_task_3
# [CALLBACK] Task executed, id: db88bbcc-9235-40ca-800c-a3ec42d1c2f2, name: test_task_1, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: 2b98b896-822d-4b2a-91cd-3b74cfdd7acd, name: test_task_4
# [CALLBACK] Task removed, id: db88bbcc-9235-40ca-800c-a3ec42d1c2f2, name: test_task_1
# [CALLBACK] Task executed, id: ea23211f-0961-49df-804e-83c93d80ce85, name: test_task_7, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: ea23211f-0961-49df-804e-83c93d80ce85, name: test_task_7
# [CALLBACK] Task executed, id: be65abf3-e1d5-403c-a180-05cb42a80fb0, name: test_task_8, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: be65abf3-e1d5-403c-a180-05cb42a80fb0, name: test_task_8
```
