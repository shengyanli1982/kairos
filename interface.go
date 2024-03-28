package kairos

import "context"

type ScheduledTaskHandleFunc = func(ctx context.Context) (data any, err error)

type ScheduledTaskCallback interface {
	OnExecuted(id, name string, data any, reason, err error)
}

type EmptyTaskCallback struct{}

func (EmptyTaskCallback) OnExecuted(id, name string, data any, reason, err error) {}

func NewEmptyTaskCallback() *EmptyTaskCallback { return &EmptyTaskCallback{} }
