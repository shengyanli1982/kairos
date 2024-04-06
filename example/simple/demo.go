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
