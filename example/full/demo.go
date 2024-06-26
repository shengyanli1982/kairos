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
	fmt.Printf("# [CALLBACK] Task executed, id: %s, name: %s, data: %v, reason: %v, err: %v\n", id, name, result, reason, err)
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
