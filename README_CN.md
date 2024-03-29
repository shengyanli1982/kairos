[English](./README.md) | 中文

<div align="center" style="position:relative;">
	<img src="assets/logo.png" alt="logo">
</div>

# 介绍

**Kairos** 源自希腊语中的“时间”一词，意为正确或适时的时刻。它是一个配备了预定义超时机制的库，用于在特定时刻执行特定任务。

我希望简化我的开发过程，并轻松设置库以在特定时间点执行任务。在长期的开发工作中，我意识到我需要围绕上下文编写大量代码，这缺乏通用性。当出现问题时，很难确定错误的来源。许多业务代码无法在不同项目之间重用。

# 为什么选择 Kairos

为了早点完成工作，有更多时间陪伴家人。出于这个目的，我将繁琐的工作任务抽象成了一个通用的函数库，方便使用。当然，我也希望能帮助到你。

**我设计的初衷是：**

1. **易于使用**：学习成本不高。
2. **高度可靠**：使用少量代码完成复杂任务，而不引入过多复杂的工具包。
3. **简单逻辑**：直接使用 Golang 的 GMP 协程模型。

# 优势

-   简单易用
-   轻量化，无需外部依赖
-   支持自定义操作的回调函数

# 安装

```bash
go get github.com/shengyanli1982/kairos
```

# 快速入门

`Kairos` 的使用非常简单，只需几行代码即可开始。

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

## 3. 示例

示例代码位于 `examples` 目录中。

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

**执行结果**

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
