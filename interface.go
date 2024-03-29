package kairos

import "time"

// WaitForContextDone 是一个只能接收的通道，用于等待上下文完成
// WaitForContextDone is a receive-only channel used to wait for context completion
type WaitForContextDone = <-chan struct{}

// TaskHandleFunc 是一个函数类型，它接收一个 WaitForContextDone 参数，并返回一个接口类型的数据和一个错误
// TaskHandleFunc is a function type that takes a WaitForContextDone parameter and returns an interface type data and an error
type TaskHandleFunc = func(done WaitForContextDone) (data interface{}, err error)

// Callback 是一个接口，定义了任务添加、执行和移除时的回调函数
// Callback is an interface that defines the callback functions when a task is added, executed, and removed
type Callback interface {
	// OnTaskAdd 是当任务被添加时的回调函数，它接收任务 id、任务名称和执行时间作为参数
	// OnTaskAdd is the callback function when a task is added, it takes the task id, task name, and execution time as parameters
	OnTaskAdd(id, name string, execAt time.Time)

	// OnTaskExecuted 是当任务被执行时的回调函数，它接收任务 id、任务名称、数据、原因和错误作为参数
	// OnTaskExecuted is the callback function when a task is executed, it takes the task id, task name, data, reason, and error as parameters
	OnTaskExecuted(id, name string, data interface{}, reason, err error)

	// OnTaskRemoved 是当任务被移除时的回调函数，它接收任务 id 和任务名称作为参数
	// OnTaskRemoved is the callback function when a task is removed, it takes the task id and task name as parameters
	OnTaskRemoved(id, name string)
}

// EmptyCallback 是一个空的回调实现，它的所有方法都是空操作
// EmptyCallback is an empty callback implementation, all of its methods are no-ops
type EmptyCallback struct{}

// OnTaskExecuted 是 EmptyCallback 的一个方法，它是一个空操作
// OnTaskExecuted is a method of EmptyCallback, it is a no-op
func (EmptyCallback) OnTaskExecuted(id, name string, data interface{}, reason, err error) {}

// OnTaskRemoved 是 EmptyCallback 的一个方法，它是一个空操作
// OnTaskRemoved is a method of EmptyCallback, it is a no-op
func (EmptyCallback) OnTaskRemoved(id, name string) {}

// OnTaskAdd 是 EmptyCallback 的一个方法，它是一个空操作
// OnTaskAdd is a method of EmptyCallback, it is a no-op
func (EmptyCallback) OnTaskAdd(id, name string, execAt time.Time) {}

// NewEmptyTaskCallback 是一个函数，它返回一个新的 EmptyCallback 实例
// NewEmptyTaskCallback is a function that returns a new instance of EmptyCallback
func NewEmptyTaskCallback() *EmptyCallback { return &EmptyCallback{} }
