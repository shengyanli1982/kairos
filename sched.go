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

	// taskCtxCache 是一个指向 cache.Cache 结构体的指针，用于存储任务的上下文信息。
	// taskCtxCache is a pointer to the cache.Cache struct, used to store the context information of tasks.
	taskCtxCache *cache.Cache

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

		// taskCtxCache 字段被设置为一个新的 Cache 结构体。
		// The taskCtxCache field is set to a new Cache struct.
		taskCtxCache: cache.NewCache(),

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
			data.(*Task).Cancel()
		})

		// 清理 taskCtxCache，取消所有任务的上下文。
		// Clean up taskCtxCache, cancel the context of all tasks.
		s.taskCtxCache.Cleanup(func(data any) {
			data.(context.CancelFunc)()
		})
	})
}

// add 是一个方法，用于向调度器添加新的任务。
// add is a method used to add new tasks to the scheduler.
func (s *Scheduler) add(ctx context.Context, cancel context.CancelFunc, name string, handleFunc TaskHandleFunc) string {
	// 创建一个新的任务。
	// Create a new task.
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
	id := task.GetMetadata().GetID()

	// 将任务添加到 taskCache 中。
	// Add the task to taskCache.
	s.taskCache.Set(id, task)

	// 将任务的取消函数添加到 taskCtxCache 中。
	// Add the cancel function of the task to taskCtxCache.
	s.taskCtxCache.Set(id, cancel)

	// 返回任务的 ID。
	// Return the ID of the task.
	return id
}

// SetAt 是一个方法，用于在指定的时间执行任务。
// SetAt is a method used to execute tasks at a specified time.
func (s *Scheduler) SetAt(name string, handleFunc TaskHandleFunc, execAt time.Time) string {
	// 创建一个新的上下文，该上下文将在指定的时间被取消。
	// Create a new context that will be cancelled at the specified time.
	ctx, cancel := context.WithDeadline(s.ctx, execAt)

	// 添加新的任务到调度器，并获取任务的 ID。
	// Add a new task to the scheduler and get the ID of the task.
	id := s.add(ctx, cancel, name, handleFunc)

	// 调用回调函数，通知任务已经被添加。
	// Call the callback function to notify that the task has been added.
	s.cfg.callback.OnTaskAdded(id, name, execAt)

	// 返回任务的 ID。
	// Return the ID of the task.
	return id
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
	// Get the task from taskCache.
	if task, ok := s.taskCache.Get(id); ok {
		// 如果任务存在，返回任务。
		// If the task exists, return the task.
		return task.(*Task)
	}

	// 如果任务不存在，返回 nil。
	// If the task does not exist, return nil.
	return nil
}

// Delete 是一个方法，用于删除指定 ID 的任务。
// Delete is a method used to delete the task with the specified ID.
func (s *Scheduler) Delete(id string) {
	taskName := ""

	// 从 taskCache 中获取任务。
	// Get the task from taskCache.
	if task, ok := s.taskCache.Get(id); ok {
		// 如果任务存在，获取任务的名称，并从 taskCache 中删除任务。
		// If the task exists, get the name of the task and delete the task from taskCache.
		taskName = task.(*Task).GetMetadata().GetName()
		s.taskCache.Delete(id)
		// 等待任务完成。
		// Wait for the task to complete.
		task.(*Task).Wait()
	}

	// 从 taskCtxCache 中获取任务的取消函数。
	// Get the cancel function of the task from taskCtxCache.
	if cancel, ok := s.taskCtxCache.Get(id); ok {
		// 如果取消函数存在，调用取消函数，并从 taskCtxCache 中删除取消函数。
		// If the cancel function exists, call the cancel function and delete the cancel function from taskCtxCache.
		cancel.(context.CancelFunc)()
		s.taskCtxCache.Delete(id)
	}

	// 如果任务名称不为空，调用回调函数，通知任务已经被删除。
	// If the task name is not empty, call the callback function to notify that the task has been deleted.
	if taskName != "" {
		s.cfg.callback.OnTaskRemoved(id, taskName)
	}
}

// GetTaskCount 是一个方法，用于获取调度器中的任务数量。
// GetTaskCount is a method used to get the number of tasks in the scheduler.
func (s *Scheduler) GetTaskCount() int {
	// 返回 taskCache 中的元素数量。
	// Return the number of elements in taskCache.
	return s.taskCache.Count()
}
