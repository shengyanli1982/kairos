package kairos

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
)

// 定义一些常见的任务错误
// Define some common task errors
var (
	// ErrorTaskCanceled 表示任务被取消
	// ErrorTaskCanceled represents the task is canceled
	ErrorTaskCanceled = errors.New("task canceled")

	// ErrorTaskTimeout 表示任务超时
	// ErrorTaskTimeout represents the task is timeout
	ErrorTaskTimeout = errors.New("task timeout")

	// ErrorTaskEarlyReturn 表示任务提前返回
	// ErrorTaskEarlyReturn represents the task returns early
	ErrorTaskEarlyReturn = errors.New("task early return")
)

// onFinishedHandleFunc 是一个函数类型，它接受一个 TaskMetadata 指针
// onFinishedHandleFunc is a function type that accepts a pointer to TaskMetadata
type onFinishedHandleFunc = func(metadata *TaskMetadata)

// onExecutedHandleFunc 是一个函数类型，它接受任务 id、name、data、reason 和 err
// onExecutedHandleFunc is a function type that accepts task id, name, data, reason and err
type onExecutedHandleFunc = func(id, name string, data any, reason, err error)

// DefaultTaskHandleFunc 是默认的任务处理函数，它返回 nil 数据和 nil 错误
// DefaultTaskHandleFunc is the default task handling function, it returns nil data and nil error
var DefaultTaskHandleFunc TaskHandleFunc = func(done WaitForContextDone) (data any, err error) { return nil, nil }

// defaultExecutedHandleFunc 是默认的执行处理函数，它不执行任何操作
// defaultExecutedHandleFunc is the default executed handling function, it does nothing
var defaultExecutedHandleFunc onExecutedHandleFunc = func(id, name string, data any, reason, err error) {}

// defaultFinishedHandleFunc 是默认的完成处理函数，它不执行任何操作
// defaultFinishedHandleFunc is the default finished handling function, it does nothing
var defaultFinishedHandleFunc onFinishedHandleFunc = func(metadata *TaskMetadata) {}

// taskPool 是一个同步池，用于存储和复用 Task 对象
// taskPool is a sync pool used to store and reuse Task objects
var taskPool = sync.Pool{New: func() interface{} { return &Task{metadata: &TaskMetadata{}} }}

// taskRefPool 是一个 sync.Pool 类型的变量，用于存储任务引用池。
// taskRefPool is a variable of type sync.Pool, used to store the task reference pool.
var taskRefPool = sync.Pool{New: func() any { return &TaskRef{parentRef: &ParentRef{}} }}

// ParentRef 结构体包含一个上下文和一个取消函数。
// The ParentRef struct contains a context and a cancel function.
type ParentRef struct {
	// ctx 是上下文对象，它可以用于传递请求范围的值、取消信号、截止时间等
	// ctx is a context object, which can be used to pass request-scoped values, cancellation signals, deadlines, etc.
	ctx context.Context

	// cancel 是一个取消函数，它可以用于取消与 ctx 关联的操作
	// cancel is a cancel function, which can be used to cancel operations associated with ctx
	cancel context.CancelFunc
}

// TaskRef 结构体包含一个父引用和一个任务。
// The TaskRef struct contains a parent reference and a task.
type TaskRef struct {
	// parentRef 是一个指向 ParentRef 的指针，它表示任务的父引用
	// parentRef is a pointer to ParentRef, which represents the parent reference of the task
	parentRef *ParentRef

	// task 是一个指向 Task 的指针，它表示任务本身
	// task is a pointer to Task, which represents the task itself
	task *Task
}

// Reset 方法重置任务引用的父引用和任务。
// The Reset method resets the parent reference and task of the task reference.
func (ref *TaskRef) Reset() {
	// 重置父引用的上下文
	// Reset the context of the parent reference
	ref.parentRef.ctx = nil

	// 重置父引用的取消函数
	// Reset the cancel function of the parent reference
	ref.parentRef.cancel = nil

	// 重置任务
	// Reset the task
	ref.task = nil
}

// TaskMetadata 结构体包含任务的 id、name 和 handleFunc
// The TaskMetadata struct contains the id, name and handleFunc of the task
type TaskMetadata struct {
	// id 是任务的唯一标识符
	// id is the unique identifier of the task
	id string

	// name 是任务的名称
	// name is the name of the task
	name string

	// handleFunc 是任务的处理函数，它定义了任务的具体执行逻辑
	// handleFunc is the handling function of the task, which defines the specific execution logic of the task
	handleFunc TaskHandleFunc
}

// GetID 方法返回任务的 id
// The GetID method returns the id of the task
func (stm *TaskMetadata) GetID() string {
	return stm.id
}

// GetName 方法返回任务的 name
// The GetName method returns the name of the task
func (stm *TaskMetadata) GetName() string {
	return stm.name
}

// GetHandleFunc 方法返回任务的 handleFunc
// The GetHandleFunc method returns the handleFunc of the task
func (stm *TaskMetadata) GetHandleFunc() TaskHandleFunc {
	return stm.handleFunc
}

// Task 结构体定义
// Definition of Task struct
type Task struct {
	// metadata 是任务的元数据，包含任务的 id、name 和 handleFunc
	// metadata is the metadata of the task, including the id, name and handleFunc of the task
	metadata *TaskMetadata

	// parentCtx 是父级上下文，用于传递上下文信息
	// parentCtx is the parent context, used to pass context information
	parentCtx context.Context

	// ctx 是任务的上下文，用于控制任务的生命周期
	// ctx is the context of the task, used to control the lifecycle of the task
	ctx context.Context

	// cancel 是一个函数，用于取消任务
	// cancel is a function used to cancel the task
	cancel context.CancelCauseFunc

	// once 用于确保任务的取消操作只执行一次
	// once is used to ensure that the cancellation operation of the task is executed only once
	once *sync.Once

	// wg 是一个 WaitGroup，用于等待任务的完成
	// wg is a WaitGroup, used to wait for the completion of the task
	wg *sync.WaitGroup

	// onFinFunc 是任务完成时的回调函数
	// onFinFunc is the callback function when the task is completed
	onFinFunc onFinishedHandleFunc

	// onExecFunc 是任务执行时的回调函数
	// onExecFunc is the callback function when the task is executed
	onExecFunc onExecutedHandleFunc
}

// NewTask 函数用于创建一个新的任务
// The NewTask function is used to create a new task
func NewTask(parentCtx context.Context, name string, handleFunc TaskHandleFunc) *Task {
	// 如果 handleFunc 为 nil，则使用默认的任务处理函数
	// If handleFunc is nil, use the default task handling function
	if handleFunc == nil {
		handleFunc = DefaultTaskHandleFunc
	}

	// 如果 parentCtx 为 nil，则使用默认的上下文
	// If parentCtx is nil, use the default context
	if parentCtx == nil {
		parentCtx = context.Background()
	}

	// 从任务池中获取一个任务
	// Get a task from the task pool
	task := taskPool.Get().(*Task)

	// 为任务生成一个新的 id
	// Generate a new id for the task
	task.metadata.id = uuid.NewString()

	// 设置任务的名称
	// Set the name of the task
	task.metadata.name = name

	// 设置任务的处理函数
	// Set the handling function of the task
	task.metadata.handleFunc = handleFunc

	// 设置任务的父级上下文
	// Set the parent context of the task
	task.parentCtx = parentCtx

	// 创建一个新的上下文和取消函数
	// Create a new context and cancel function
	task.ctx, task.cancel = context.WithCancelCause(parentCtx)

	// 创建一个新的 WaitGroup
	// Create a new WaitGroup
	task.wg = &sync.WaitGroup{}

	// 创建一个新的 Once
	// Create a new Once
	task.once = &sync.Once{}

	// 增加 WaitGroup 的计数
	// Increase the count of WaitGroup
	task.wg.Add(1)

	// 在一个新的 goroutine 中执行任务
	// Execute the task in a new goroutine
	go task.executor()

	// 返回任务
	// Return the task
	return task
}

// executor 方法用于执行任务
// The executor method is used to execute the task
func (t *Task) executor() {
	// 创建一个定时器，每 500 毫秒触发一次
	// Create a timer that triggers every 500 milliseconds
	taskTrigger := time.NewTicker(time.Millisecond * 500)

	// 使用 defer 语句确保在函数结束时执行一些操作
	// Use the defer statement to ensure that some operations are performed at the end of the function
	defer func() {
		// 停止定时器
		// Stop the timer
		taskTrigger.Stop()

		// 减少 WaitGroup 的计数
		// Decrease the count of WaitGroup
		t.wg.Done()

		// 将任务放回任务池
		// Put the task back into the task pool
		taskPool.Put(t)

		// 调用 onFinFunc 回调函数，传入任务的元数据
		// Call the onFinFunc callback function, passing in the metadata of the task
		t.onFinFunc(t.metadata)
	}()

	// 使用 for 循环和 select 语句来监听和处理事件
	// Use a for loop and select statement to listen for and handle events
	for {
		select {
		// 当任务的上下文被取消时
		// When the context of the task is canceled
		case <-t.ctx.Done():
			// 获取取消的原因
			// Get the reason for the cancellation
			reason := context.Cause(t.ctx)

			// 根据取消的原因来处理任务
			// Handle the task based on the reason for the cancellation
			switch reason {
			// 如果任务被取消
			// If the task is canceled
			case context.Canceled:
				// 调用 onExecFunc 回调函数，传入任务 id、任务名称、nil 结果、任务取消错误和 nil 错误
				// Call the onExecFunc callback function, passing in the task id, task name, nil result, task cancellation error, and nil error
				t.onExecFunc(t.metadata.id, t.metadata.name, nil, ErrorTaskCanceled, nil)

			// 如果任务超时
			// If the task is timeout
			case context.DeadlineExceeded:
				// 调用任务的处理函数，获取结果和错误
				// Call the task's handling function to get the result and error
				result, err := t.metadata.handleFunc(t.ctx.Done())

				// 调用 onExecFunc 回调函数，传入任务 id、任务名称、结果、任务超时错误和错误
				// Call the onExecFunc callback function, passing in the task id, task name, result, task timeout error, and error
				t.onExecFunc(t.metadata.id, t.metadata.name, result, ErrorTaskTimeout, err)

			// 如果任务提前返回
			// If the task returns early
			case ErrorTaskEarlyReturn:
				// 调用任务的处理函数，获取结果和错误
				// Call the task's handling function to get the result and error
				result, err := t.metadata.handleFunc(t.ctx.Done())

				// 调用 onExecFunc 回调函数，传入任务 id、任务名称、结果、任务提前返回错误和错误
				// Call the onExecFunc callback function, passing in the task id, task name, result, task early return error, and error
				t.onExecFunc(t.metadata.id, t.metadata.name, result, ErrorTaskEarlyReturn, err)
			}

			// 取消任务
			// Cancel the task
			t.cancel(reason)

			// 结束函数
			// End the function
			return

		// 当定时器触发时
		// When the timer triggers
		case <-taskTrigger.C:
			// 让出 CPU 时间片，让其他 goroutine 有机会运行
			// Yield the CPU time slice to give other goroutines a chance to run
			runtime.Gosched()
		}
	}
}

// EarlyReturn 方法用于提前返回任务
// The EarlyReturn method is used to return the task early
func (t *Task) EarlyReturn() {
	// 如果 once 不为 nil
	// If once is not nil
	if t.once != nil {
		// 使用 once.Do 方法确保 cancel 方法只被调用一次
		// Use the once.Do method to ensure that the cancel method is called only once
		t.once.Do(func() {
			// 调用 cancel 方法，传入 ErrorTaskEarlyReturn 错误
			// Call the cancel method, passing in the ErrorTaskEarlyReturn error
			t.cancel(ErrorTaskEarlyReturn)
		})
	}
}

// Cancel 方法用于取消任务
// The Cancel method is used to cancel the task
func (t *Task) Cancel() {
	// 如果 once 不为 nil
	// If once is not nil
	if t.once != nil {
		// 使用 once.Do 方法确保 cancel 方法只被调用一次
		// Use the once.Do method to ensure that the cancel method is called only once
		t.once.Do(func() {
			// 调用 cancel 方法，传入 context.Canceled 错误
			// Call the cancel method, passing in the context.Canceled error
			t.cancel(context.Canceled)
		})
	}
}

// GetMetadata 方法用于获取任务的元数据
// The GetMetadata method is used to get the metadata of the task
func (t *Task) GetMetadata() *TaskMetadata {
	// 返回任务的元数据
	// Return the metadata of the task
	return t.metadata
}

// Wait 方法用于等待任务完成
// The Wait method is used to wait for the task to complete
func (t *Task) Wait() {
	// 调用 WaitGroup 的 Wait 方法
	// Call the Wait method of WaitGroup
	t.wg.Wait()
}

// onFinished 方法用于设置任务完成时的回调函数
// The onFinished method is used to set the callback function when the task is completed
func (t *Task) onFinished(fn onFinishedHandleFunc) *Task {
	// 如果 fn 为 nil
	// If fn is nil
	if fn == nil {
		// 使用默认的完成处理函数
		// Use the default finished handling function
		fn = defaultFinishedHandleFunc
	}

	// 设置 onFinFunc
	// Set onFinFunc
	t.onFinFunc = fn

	// 返回任务
	// Return the task
	return t
}

// onExecuted 方法用于设置任务执行时的回调函数
// The onExecuted method is used to set the callback function when the task is executed
func (t *Task) onExecuted(fn onExecutedHandleFunc) *Task {
	// 如果 fn 为 nil
	// If fn is nil
	if fn == nil {
		// 使用默认的执行处理函数
		// Use the default executed handling function
		fn = defaultExecutedHandleFunc
	}

	// 设置 onExecFunc
	// Set onExecFunc
	t.onExecFunc = fn

	// 返回任务
	// Return the task
	return t
}
