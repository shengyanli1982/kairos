package kairos

import (
	"context"
	"sync"
	"time"

	"github.com/shengyanli1982/kairos/internal/cache"
)

// Scheduler 结构体定义了一个调度器，它包含了一些用于任务调度的关键字段。
// The Scheduler struct defines a scheduler, which contains some key fields for task scheduling.
type Scheduler struct {
	// cfg 是一个指向 Config 结构体的指针，用于存储调度器的配置信息。
	// cfg is a pointer to the Config struct, used to store the configuration information of the scheduler.
	cfg *Config

	// taskCache 是一个指向 cache.Cache 结构体的指针，用于存储已经调度的任务。
	// taskCache is a pointer to the cache.Cache struct, used to store the scheduled tasks.
	taskCache *cache.Cache

	// uniqCache 是一个指向 cache.Cache 结构体的指针，用于存储唯一的任务。
	// uniqCache is a pointer to the cache.Cache struct, used to store unique tasks.
	uniqCache *cache.Cache

	// ctx 是一个 context.Context 类型的变量，用于存储调度器的上下文信息。
	// ctx is a variable of type context.Context, used to store the context information of the scheduler.
	ctx context.Context

	// cancel 是一个 context.CancelFunc 类型的函数，用于取消调度器的上下文。
	// cancel is a function of type context.CancelFunc, used to cancel the context of the scheduler.
	cancel context.CancelFunc

	// once 是一个 sync.Once 类型的变量，用于确保某个操作只执行一次。
	// once is a variable of type sync.Once, used to ensure that a certain operation is performed only once.
	once sync.Once
}

// New 是一个函数，接收一个指向 Config 结构体的指针作为参数，返回一个新的 Scheduler 结构体指针。
// New is a function that takes a pointer to a Config struct as a parameter and returns a new pointer to a Scheduler struct.
func New(conf *Config) *Scheduler {
	// 首先，我们检查传入的配置是否有效，如果无效，将返回一个默认的配置。
	// First, we check if the passed configuration is valid, if not, a default configuration will be returned.
	conf = isConfigValid(conf)

	// 然后，我们创建一个新的 Scheduler 结构体，并初始化它的字段。
	// Then, we create a new Scheduler struct and initialize its fields.
	s := &Scheduler{
		// cfg 字段被设置为传入的配置。
		// The cfg field is set to the passed configuration.
		cfg: conf,

		// taskCache 字段被设置为一个新的 Cache 结构体。
		// The taskCache field is set to a new Cache struct.
		taskCache: cache.NewCache(),

		// uniqCache 字段被设置为一个新的 Cache 结构体。
		// The uniqCache field is set to a new Cache struct.
		uniqCache: cache.NewCache(),

		// once 字段被设置为一个新的 Once 结构体。
		// The once field is set to a new Once struct.
		once: sync.Once{},
	}

	// ctx 和 cancel 字段被设置为一个新的带取消功能的上下文。
	// The ctx and cancel fields are set to a new context with cancellation.
	s.ctx, s.cancel = context.WithCancel(context.Background())

	// 最后，我们返回新创建的 Scheduler 结构体的指针。
	// Finally, we return the pointer to the newly created Scheduler struct.
	return s
}

// Stop 是一个方法，用于停止调度器的所有任务。
// Stop is a method used to stop all tasks of the scheduler.
func (s *Scheduler) Stop() {
	// 使用 sync.Once 确保停止操作只执行一次。
	// Use sync.Once to ensure that the stop operation is performed only once.
	s.once.Do(func() {
		// 调用 cancel 函数来取消调度器的上下文，从而停止所有任务。
		// Call the cancel function to cancel the context of the scheduler, thereby stopping all tasks.
		s.cancel()

		// 清理 taskCache，取消所有已经调度的任务。
		// Clean up taskCache, cancel all scheduled tasks.
		s.taskCache.Cleanup(func(data any) {
			// 将数据转换为任务引用。
			// Convert the data to a task reference.
			taskRef := data.(*TaskRef)

			// 取消任务。
			// Cancel the task.
			taskRef.task.Cancel()

			// 取消父引用的上下文。
			// Cancel the context of the parent reference.
			taskRef.parentRef.cancel()

			// 等待任务完成。
			// Wait for the task to complete.
			taskRef.task.Wait()

			// 重置任务引用。
			// Reset the task reference.
			taskRef.Reset()

			// 将任务引用放回到任务引用池中。
			// Put the task reference back into the task reference pool.
			taskRefPool.Put(taskRef)
		})

		// 如果调度器的配置中 uniqTask 为 true
		// If uniqTask in the scheduler's configuration is true
		if s.cfg.unique {
			// 调用 uniqCache 的 Cleanup 方法，清理其中的所有任务
			// Call the Cleanup method of uniqCache to clean up all the tasks in it
			s.uniqCache.Cleanup(func(value any) {})
		}
	})
}

// add 是一个方法，用于向调度器添加新的任务。
// add is a method used to add new tasks to the scheduler.
func (s *Scheduler) add(ctx context.Context, cancel context.CancelFunc, name string, handleFunc TaskHandleFunc) string {
	// 如果调度器的配置中 uniqTask 为 true
	// If uniqTask in the scheduler's configuration is true
	if s.cfg.unique {
		// 从 uniqCache 中获取任务
		// Get the task from uniqCache
		if data, ok := s.uniqCache.Get(name); ok {
			// 获取任务的 ID
			// Get the ID of the task
			taskID := data.(string)

			// 调用回调函数，通知任务已经存在。
			// Call the callback function to notify that the task already exists.
			s.cfg.callback.OnTaskDuplicated(taskID, name)

			// 返回任务的 ID。
			// Return the ID of the task.
			return taskID
		}
	}

	// 创建一个新的任务，并设置任务执行后和任务完成后的回调函数。
	// Create a new task, and set the callback functions after the task is executed and after the task is finished.
	task := NewTask(ctx, name, handleFunc).
		// 设置任务执行后的回调函数。
		// Set the callback function after the task is executed.
		onExecuted(s.cfg.callback.OnTaskExecuted).

		// 设置任务完成后的回调函数。
		// Set the callback function after the task is finished.
		onFinished(func(metadata *TaskMetadata) {
			// 任务完成后，从调度器中删除该任务。
			// After the task is finished, delete the task from the scheduler.
			s.Delete(metadata.GetID())
		})

	// 获取任务的 ID。
	// Get the ID of the task.
	taskID := task.GetMetadata().GetID()

	// 如果调度器的配置中 unique 为 true
	// If unique in the scheduler's configuration is true
	if s.cfg.unique {
		// 在 uniqCache 中设置该任务的 ID
		// Set the ID of the task in uniqCache
		s.uniqCache.Set(name, taskID)
	}

	// 从任务引用池中获取一个任务引用。
	// Get a task reference from the task reference pool.
	taskRef := taskRefPool.Get().(*TaskRef)

	// 设置任务引用的父引用的上下文。
	// Set the context of the parent reference of the task reference.
	taskRef.parentRef.ctx = ctx

	// 设置任务引用的父引用的取消函数。
	// Set the cancel function of the parent reference of the task reference.
	taskRef.parentRef.cancel = cancel

	// 设置任务引用的任务。
	// Set the task of the task reference.
	taskRef.task = task

	// 在任务缓存中设置任务引用。
	// Set the task reference in the task cache.
	s.taskCache.Set(taskID, taskRef)

	// 返回任务的 ID。
	// Return the ID of the task.
	return taskID
}

// SetAt 是一个方法，用于在指定时间执行任务。
// SetAt is a method used to execute tasks at a specified time.
func (s *Scheduler) SetAt(name string, handleFunc TaskHandleFunc, execAt time.Time) string {
	// 创建一个新的上下文，该上下文将在指定时间被取消。
	// Create a new context that will be cancelled at the specified time.
	ctx, cancel := context.WithDeadline(s.ctx, execAt)

	// 添加一个新的任务到调度器，并获取任务的 ID。
	// Add a new task to the scheduler and get the ID of the task.
	taskID := s.add(ctx, cancel, name, handleFunc)

	// 调用回调函数，通知任务已被添加。
	// Call the callback function to notify that the task has been added.
	s.cfg.callback.OnTaskAdded(taskID, name, execAt)

	// 返回任务的 ID。
	// Return the ID of the task.
	return taskID
}

// Set 是一个方法，用于在指定的延迟后执行任务。
// Set is a method used to execute tasks after a specified delay.
func (s *Scheduler) Set(name string, handleFunc TaskHandleFunc, delay time.Duration) string {
	// 调用 SetAt 方法，将当前时间加上指定的延迟作为执行时间。
	// Call the SetAt method, adding the specified delay to the current time as the execution time.
	return s.SetAt(name, handleFunc, time.Now().Add(delay))
}

// Get 是一个方法，用于获取指定 ID 的任务。
// Get is a method used to get the task with the specified ID.
func (s *Scheduler) Get(id string) *Task {
	// 从 taskCache 中获取任务。
	// Get the data from taskCache.
	if data, ok := s.taskCache.Get(id); ok {
		// 如果任务存在，返回任务。
		// If the task exists, return the task.
		return data.(*TaskRef).task
	}

	// 如果任务不存在，返回 nil。
	// If the task does not exist, return nil.
	return nil
}

// Delete 是一个方法，用于删除指定 ID 的任务。
// Delete is a method used to delete the task with the specified ID.
func (s *Scheduler) Delete(id string) {
	// 从 taskCache 中获取任务。
	// Get the task from taskCache.
	if data, ok := s.taskCache.Get(id); ok {
		// 如果获取成功，将数据转换为 Task 类型。
		// If the retrieval is successful, convert the data to the Task type.
		taskRef := data.(*TaskRef)

		// 调用任务的 Cancel 方法来取消任务。
		// Call the Cancel method of the task to cancel the task.
		taskRef.task.Cancel()

		// 调用父引用的 cancel 函数来取消上下文。
		// Call the cancel function of the parent reference to cancel the context.
		taskRef.parentRef.cancel()

		// 获取任务的名称。
		// Get the name of the task.
		taskName := taskRef.task.GetMetadata().GetName()

		// 从任务缓存中删除这个任务。
		// Delete this task from the task cache.
		s.taskCache.Delete(id)

		// 调用任务的 Wait 方法来等待任务完成。
		// Call the Wait method of the task to wait for the task to complete.
		taskRef.task.Wait()

		// 重置任务引用。
		// Reset the task reference.
		taskRef.Reset()

		// 将任务引用放回任务引用池。
		// Put the task reference back into the task reference pool.
		taskRefPool.Put(taskRef)

		// 调用回调函数，通知任务已经被删除。
		// Call the callback function to notify that the task has been deleted.
		s.cfg.callback.OnTaskRemoved(id, taskName)

		// 如果调度器的配置中 uniqTask 为 true
		// If uniqTask in the scheduler's configuration is true
		if s.cfg.unique {
			// 从 uniqCache 中删除指定名称的任务
			// Delete the task with the specified name from uniqCache
			s.uniqCache.Delete(taskName)
		}
	}
}

// Count 是一个方法，用于获取调度器中的任务数量。
// Count is a method used to get the number of tasks in the scheduler.
func (s *Scheduler) Count() int {
	// 返回 taskCache 中的元素数量。
	// Return the number of elements in taskCache.
	return s.taskCache.Count()
}
