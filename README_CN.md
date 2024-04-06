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

	// OnTaskExecuted 是当任务被执行时的回调函数，它接收任务 id、任务名称、数据、原因和错误作为参数
	// OnTaskExecuted is the callback function when a task is executed, it takes the task id, task name, data, reason, and error as parameters
	OnTaskExecuted(id, name string, data interface{}, reason, err error)

	// OnTaskRemoved 是当任务被移除时的回调函数，它接收任务 id 和任务名称作为参数
	// OnTaskRemoved is the callback function when a task is removed, it takes the task id and task name as parameters
	OnTaskRemoved(id, name string)

	// OnTaskDuplicated 是当任务重复时的回调函数，它接收任务 id 和任务名称作为参数
	// OnTaskDuplicated is the callback function when a task is duplicated, it takes the task id and task name as parameters
	OnTaskDuplicated(id, name string)
}
```

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
> 如果您希望确保 `Scheduler` 中的任务是唯一的，可以使用 `WithDisableDuplicated` 选项来禁用重复的任务。
>
> 默认情况下，`Scheduler` 允许重复的任务（`WithDisableDuplicated` 设置为 `false`）。这意味着您可以向 `Scheduler` 添加具有相同名称的多个任务。
>
> 如果希望 `Scheduler` 使用任务名称作为标识符，请确保在 `WithDisableDuplicated` 设置为 `true` 时为每个任务使用不同的名称。
>
> 当 `WithDisableDuplicated` 设置为 `true` 时，`Set` 和 `SetAt` 方法将返回正在运行的任务的 `id`。

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
	config := ks.NewConfig().WithCallback(&demoSchedCallback{}).WithDisableDuplicated(true)

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

	// 添加一个名字存在的任务到调度器，等待触发回调。
	// Add a task with an existing name to the scheduler and wait for the callback to be triggered.
	taskID := scheduler.Set("test_task_9", func(done ks.WaitForContextDone) (result any, err error) {
		// 当任务完成时，返回任务的名称 (这步不是必须，只是为了表示函数内部可以接收外部的 ctx 信号)。
		// When the task is done, return the name of the task (this step is not necessary, just to indicate that the function can receive the ctx signal from the outside).
		for range done {
			return "test_task_9", nil
		}

		// 如果任务没有完成，返回 nil。
		// If the task is not done, return nil.
		return nil, nil

		// 设置任务的延迟时间为 200 毫秒。
		// Set the delay time of the task to 200 milliseconds.
	}, time.Millisecond*200)

	// 打印任务添加的信息。
	// Print the information when the task is added.
	fmt.Printf("%% [MAIN] The duplicate task can be retrieved, id: %s, name: %s\n", taskID, "test_task_1")

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
# [CALLBACK] Task added, id: ec1c1cc9-982b-4bad-b818-3be9f45120a6, name: test_task_0, execAt: 2024-04-06 16:40:06.769283 +0800 CST m=+0.200283861
% [MAIN] Task 0 can be retrieved, id: ec1c1cc9-982b-4bad-b818-3be9f45120a6, name: test_task_0
# [CALLBACK] Task added, id: eb63e1f1-4c7f-4664-9b37-6b0d4c087e22, name: test_task_1, execAt: 2024-04-06 16:40:06.769554 +0800 CST m=+0.200554984
% [MAIN] Task 1 can be retrieved, id: eb63e1f1-4c7f-4664-9b37-6b0d4c087e22, name: test_task_1
# [CALLBACK] Task added, id: 39219bdf-06ee-45bc-bc41-e4ee76b64987, name: test_task_2, execAt: 2024-04-06 16:40:06.769577 +0800 CST m=+0.200577535
% [MAIN] Task 2 can be retrieved, id: 39219bdf-06ee-45bc-bc41-e4ee76b64987, name: test_task_2
# [CALLBACK] Task added, id: dd0e2ce3-8e4e-49e5-a543-c0cd86202880, name: test_task_3, execAt: 2024-04-06 16:40:06.76959 +0800 CST m=+0.200590233
% [MAIN] Task 3 can be retrieved, id: dd0e2ce3-8e4e-49e5-a543-c0cd86202880, name: test_task_3
# [CALLBACK] Task added, id: 8d44a0a1-f6af-4b86-a202-bff58a25db0d, name: test_task_4, execAt: 2024-04-06 16:40:06.769609 +0800 CST m=+0.200609710
% [MAIN] Task 4 can be retrieved, id: 8d44a0a1-f6af-4b86-a202-bff58a25db0d, name: test_task_4
# [CALLBACK] Task added, id: f8482885-1d4b-45c2-94f5-23dfddb87d24, name: test_task_5, execAt: 2024-04-06 16:40:06.769671 +0800 CST m=+0.200671867
% [MAIN] Task 5 can be retrieved, id: f8482885-1d4b-45c2-94f5-23dfddb87d24, name: test_task_5
# [CALLBACK] Task added, id: 1b6a76bd-949c-4467-ab38-c1e9aa9576a6, name: test_task_6, execAt: 2024-04-06 16:40:06.769699 +0800 CST m=+0.200699429
% [MAIN] Task 6 can be retrieved, id: 1b6a76bd-949c-4467-ab38-c1e9aa9576a6, name: test_task_6
# [CALLBACK] Task added, id: 1aa633af-b22b-4ebf-a5bb-4a1410ccaa51, name: test_task_7, execAt: 2024-04-06 16:40:06.769712 +0800 CST m=+0.200712896
% [MAIN] Task 7 can be retrieved, id: 1aa633af-b22b-4ebf-a5bb-4a1410ccaa51, name: test_task_7
# [CALLBACK] Task added, id: cea3423d-e468-4bd9-bca8-197b1e3ce453, name: test_task_8, execAt: 2024-04-06 16:40:06.769733 +0800 CST m=+0.200733067
% [MAIN] Task 8 can be retrieved, id: cea3423d-e468-4bd9-bca8-197b1e3ce453, name: test_task_8
# [CALLBACK] Task added, id: d5442471-2730-4c61-936d-9f7fd89014de, name: test_task_9, execAt: 2024-04-06 16:40:06.769862 +0800 CST m=+0.200862581
% [MAIN] Task 9 can be retrieved, id: d5442471-2730-4c61-936d-9f7fd89014de, name: test_task_9
>>> [CALLBACK] Task duplicated , id: d5442471-2730-4c61-936d-9f7fd89014de, name: test_task_9
# [CALLBACK] Task added, id: d5442471-2730-4c61-936d-9f7fd89014de, name: test_task_9, execAt: 2024-04-06 16:40:06.769883 +0800 CST m=+0.200883229
% [MAIN] The duplicate task can be retrieved, id: d5442471-2730-4c61-936d-9f7fd89014de, name: test_task_1
# [CALLBACK] Task executed, id: 1aa633af-b22b-4ebf-a5bb-4a1410ccaa51, name: test_task_7, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task executed, id: cea3423d-e468-4bd9-bca8-197b1e3ce453, name: test_task_8, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task executed, id: dd0e2ce3-8e4e-49e5-a543-c0cd86202880, name: test_task_3, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task executed, id: d5442471-2730-4c61-936d-9f7fd89014de, name: test_task_9, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: d5442471-2730-4c61-936d-9f7fd89014de, name: test_task_9
# [CALLBACK] Task executed, id: f8482885-1d4b-45c2-94f5-23dfddb87d24, name: test_task_5, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: f8482885-1d4b-45c2-94f5-23dfddb87d24, name: test_task_5
# [CALLBACK] Task removed, id: cea3423d-e468-4bd9-bca8-197b1e3ce453, name: test_task_8
# [CALLBACK] Task executed, id: 8d44a0a1-f6af-4b86-a202-bff58a25db0d, name: test_task_4, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: 8d44a0a1-f6af-4b86-a202-bff58a25db0d, name: test_task_4
# [CALLBACK] Task removed, id: dd0e2ce3-8e4e-49e5-a543-c0cd86202880, name: test_task_3
# [CALLBACK] Task executed, id: 39219bdf-06ee-45bc-bc41-e4ee76b64987, name: test_task_2, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task executed, id: 1b6a76bd-949c-4467-ab38-c1e9aa9576a6, name: test_task_6, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: 1b6a76bd-949c-4467-ab38-c1e9aa9576a6, name: test_task_6
# [CALLBACK] Task executed, id: ec1c1cc9-982b-4bad-b818-3be9f45120a6, name: test_task_0, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task executed, id: eb63e1f1-4c7f-4664-9b37-6b0d4c087e22, name: test_task_1, data: <nil>, reason: task timeout, err: <nil>
# [CALLBACK] Task removed, id: 39219bdf-06ee-45bc-bc41-e4ee76b64987, name: test_task_2
# [CALLBACK] Task removed, id: ec1c1cc9-982b-4bad-b818-3be9f45120a6, name: test_task_0
# [CALLBACK] Task removed, id: eb63e1f1-4c7f-4664-9b37-6b0d4c087e22, name: test_task_1
# [CALLBACK] Task removed, id: 1aa633af-b22b-4ebf-a5bb-4a1410ccaa51, name: test_task_7
```
