package kairos

type WaitForCtxDone = <-chan struct{}

type ScheduledTaskHandleFunc = func(done WaitForCtxDone) (data any, err error)

type ScheduledTaskCallback interface {
	OnExecuted(id, name string, data any, reason, err error)
}

type EmptyTaskCallback struct{}

func (EmptyTaskCallback) OnExecuted(id, name string, data any, reason, err error) {}

func NewEmptyTaskCallback() *EmptyTaskCallback { return &EmptyTaskCallback{} }
