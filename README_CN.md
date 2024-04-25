[English](./README.md) | 中文

<div align="center" style="position:relative;">
	<img src="assets/logo.png" alt="logo">
</div>

[![Go Report Card](https://goreportcard.com/badge/github.com/shengyanli1982/kairos)](https://goreportcard.com/report/github.com/shengyanli1982/kairos)
[![Build Status](https://github.com/shengyanli1982/kairos/actions/workflows/test.yaml/badge.svg)](https://github.com/shengyanli1982/kairos/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/shengyanli1982/kairos.svg)](https://pkg.go.dev/github.com/shengyanli1982/kairos)

# 简介

**Kairos** 源自希腊语中的“时间”一词，意为正确或适时的时刻。它是一个带有预定义超时机制的库，用于在特定时刻执行特定任务。

我希望简化我的开发过程，并轻松设置库以在特定时间点执行任务。在长期的开发工作中，我意识到我需要围绕上下文编写大量代码，这缺乏通用性。当出现问题时，很难确定错误的源头。许多业务代码无法在不同项目之间重用。

# 为什么选择 Kairos

为了提前完成工作并有更多时间陪伴家人。出于这个目的，我将繁琐的工作任务抽象成一个通用的函数库，方便使用。当然，我也希望能帮助到你。

**我设计的初衷是：**

1. **易于使用**：不需要花费太多时间学习。
2. **高度可靠**：使用少量的代码完成复杂的任务，而不引入过多复杂的工具包。
3. **简单逻辑**：直接使用 Golang 的 GMP 协程模型。

# 优势

-   简单易用
-   轻量化，无外部依赖
-   支持自定义操作的回调函数

# 安装

```bash
go get github.com/shengyanli1982/kairos
```

# 快速入门

`Kairos` 使用非常简单，只需几行代码即可开始使用。

## 1. 配置

`Kairos` 有一个配置对象，可以用于注册回调函数。配置对象具有以下字段：

-   `WithCallback`：为任务注册回调函数。

    ```go
    // Callback 是一个接口，定义了任务添加、执行和移除时的回调函数
    // Callback is an interface that defines the callback functions when a task is added, executed, and removed
    type Callback interface {
    	// OnTaskAdded 是当任务被添加时的回调函数，它接收任务 id、任务名称和执行时间作为参数
    	// OnTaskAdded is the callback function when a task is added, it takes the task id, task name, and execution time as parameters
    	OnTaskAdded(id, name string, execAt time.Time)

    	// OnTaskExecuted 是当任务被执行时的回调函数，它接收任务 id、任务名称、任务结果、原因和错误作为参数
    	// OnTaskExecuted is the callback function when a task is executed, it takes the task id, task name, task result, reason, and error as parameters
    	OnTaskExecuted(id, name string, data interface{}, reason, err error)

    	// OnTaskRemoved 是当任务被移除时的回调函数，它接收任务 id 和任务名称作为参数
    	// OnTaskRemoved is the callback function when a task is removed, it takes the task id and task name as parameters
    	OnTaskRemoved(id, name string)

    	// OnTaskDuplicated 是当任务重复时的回调函数，它接收任务 id 和任务名称作为参数
    	// OnTaskDuplicated is the callback function when a task is duplicated, it takes the task id and task name as parameters
    	OnTaskDuplicated(id, name string)
    }
    ```

-   `WithUniqued`: 禁用重复任务。当设置为 `true` 时，`Scheduler` 将不允许具有相同名称的任务。

## 2. 方法

`Kairos` 提供以下方法：

-   `New`：创建一个新的 `Scheduler` 对象。`Scheduler` 对象用于管理任务。
-   `Stop`：停止 `Scheduler`。如果 `Scheduler` 对象被停止，所有任务将被停止并移除。
-   `Set`：向 `Scheduler` 添加一个任务。`Set` 方法接受任务的 `name`、执行任务的延迟时间 `delay`（time.Duration）和任务的处理函数 `handleFunc` 作为参数。
-   `SetAt`：在特定时间向 `Scheduler` 添加一个任务。`SetAt` 方法接受任务的 `name`、执行任务的时间 `execAt`（time.Time）和任务的处理函数 `handleFunc` 作为参数。
-   `Get`：通过任务的 `id` 从 `Scheduler` 获取任务。
-   `Delete`：通过任务的 `id` 从 `Scheduler` 删除任务。
-   `Count`: 获取 `Scheduler` 中任务的数量。

> [!TIP]
>
> 如果您希望确保 `Scheduler` 中的任务是唯一的，可以使用 `WithUniqued` 选项来禁用重复的任务。
>
> 默认情况下，`Scheduler` 允许重复的任务（`WithUniqued` 设置为 `false`）。这意味着您可以向 `Scheduler` 添加具有相同名称的多个任务。
>
> 如果希望 `Scheduler` 使用任务名称作为标识符，请确保在 `WithUniqued` 设置为 `true` 时为每个任务使用不同的名称。
>
> 当 `WithUniqued` 设置为 `true` 时，`Set` 和 `SetAt` 方法将返回正在运行的任务的 `id`。
>
> 如果您想要访问由 `Set` 和 `SetAt` 设置的自定义处理函数的结果值，您可以利用 `Callback` 中的 `OnTaskExecuted` 方法。该方法有一个 `result` 参数，表示自定义处理函数返回的结果值。

## 3. 任务

`Task` 是 `Kairos` 中的一个关键概念，它允许在指定的时间执行特定任务。`Task` 对象提供以下方法：

-   `GetMetadata`：获取任务的元数据，包括获取任务信息的方法。
    1.  `GetID`：获取任务的 `id`。
    2.  `GetName`：获取任务的名称。
    3.  `GetHandleFunc`：获取任务的处理函数。
-   `EarlyReturn`：手动停止任务执行并提前返回，无需等待超时或取消信号。它会调用 `handleFunc`。
-   `Cancel`：手动停止任务执行并立即返回，不执行 `handleFunc`。
-   `Wait`：等待任务完成，阻塞当前 goroutine 直到任务完成。

> [!NOTE]
>
> `Wait` 方法是一个阻塞方法。建议在单独的 `goroutine` 中使用它，以避免阻塞主 `goroutine`。
>
> 如果需要获取任务，可以使用 `Get` 方法通过任务的 `ID` 获取任务。一旦获取到任务，可以使用 `Wait` 方法等待任务完成。
>
> 不需要手动停止任务。任务执行后，其状态会自动在 `Scheduler` 中清除。如果需要主动删除任务，可以使用 `Delete` 方法删除相应的任务（仅在必要时才应该这样做，不推荐）。

## 4. 示例

示例代码位于 `examples` 目录中。

### 4.1 精简实例

只需几行代码，您就可以轻松创建和运行任务，并获取其信息。

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

**执行结果**

```bash
$ go run demo.go
% [MAIN] Task can be retrieved, id: 045a6ecb-5d4d-4393-84a3-679c943bd6f7, name: test_task
```

### 4.2 完整示例

这个示例展示了使用回调函数来处理任务的执行、添加和删除。

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

**执行结果**

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
