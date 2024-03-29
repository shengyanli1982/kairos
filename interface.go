package kairos

import "time"

type WaitForContextDone = <-chan struct{}

type TaskHandleFunc = func(done WaitForContextDone) (data interface{}, err error)

type Callback interface {
	OnTaskAdd(id, name string, execAt time.Time)
	OnTaskExecuted(id, name string, data interface{}, reason, err error)
	OnTaskRemoved(id, name string)
}

type EmptyCallback struct{}

func (EmptyCallback) OnTaskExecuted(id, name string, data interface{}, reason, err error) {}

func (EmptyCallback) OnTaskRemoved(id, name string) {}

func (EmptyCallback) OnTaskAdd(id, name string, execAt time.Time) {}

func NewEmptyTaskCallback() *EmptyCallback { return &EmptyCallback{} }
