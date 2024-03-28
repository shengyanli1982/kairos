package kairos

import "context"

type TaskHandleFunc = func(ctx context.Context) (data any, err error)

type TaskCallback interface {
	OnExecuted(id, name string, data any, reason, err error)
}

type EmptyTaskCallback struct{}

func (EmptyTaskCallback) OnExecuted(id, name string, data any, reason, err error) {}

func NewEmptyTaskCallback() *EmptyTaskCallback { return &EmptyTaskCallback{} }
