package kairos

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/google/uuid"
)

var (
	ErrorTaskCanceled    = errors.New("task canceled")
	ErrorTaskTimeout     = errors.New("task timeout")
	ErrorTaskEarlyReturn = errors.New("task early return")
)

var DefaultTaskHandleFunc TaskHandleFunc = func(ctx context.Context) (data any, err error) { return nil, nil }

var taskPool = sync.Pool{New: func() interface{} { return &Task{once: new(sync.Once)} }}

type Task struct {
	id         string
	name       string
	callback   TaskCallback
	handleFunc TaskHandleFunc
	parentCtx  context.Context
	ctx        context.Context
	cancel     context.CancelCauseFunc
	once       *sync.Once
	wg         sync.WaitGroup
}

func NewTask(parentCtx context.Context, name string, handleFunc TaskHandleFunc, callback TaskCallback) *Task {
	if handleFunc == nil {
		handleFunc = DefaultTaskHandleFunc
	}

	if parentCtx == nil {
		parentCtx = context.Background()
	}

	if callback == nil {
		callback = NewEmptyTaskCallback()
	}

	task := taskPool.Get().(*Task)

	task.id = uuid.NewString()
	task.name = name
	task.callback = callback
	task.handleFunc = handleFunc
	task.parentCtx = parentCtx
	task.ctx, task.cancel = context.WithCancelCause(parentCtx)

	task.wg.Add(1)
	go task.executor()

	return task
}

func (t *Task) Stop() {
	t.once.Do(func() {

		t.cancel(ErrorTaskEarlyReturn)

		t.wg.Wait()

		taskPool.Put(t)
	})

	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&t.once)), unsafe.Pointer(new(sync.Once)))
}

func (t *Task) executor() {

	taskTrigger := time.NewTicker(time.Millisecond * 500)

	defer func() {
		taskTrigger.Stop()
		t.wg.Done()
	}()

	for {
		select {

		case <-t.ctx.Done():

			switch context.Cause(t.ctx) {

			case context.Canceled:
				t.callback.OnExecuted(t.id, t.name, nil, ErrorTaskCanceled, nil)
				t.cancel(context.Canceled)

			case context.DeadlineExceeded:
				result, err := t.handleFunc(t.ctx)
				t.callback.OnExecuted(t.id, t.name, result, ErrorTaskTimeout, err)

			case ErrorTaskEarlyReturn:
				result, err := t.handleFunc(t.ctx)
				t.callback.OnExecuted(t.id, t.name, result, ErrorTaskEarlyReturn, err)
			}

			return

		case <-taskTrigger.C:
			runtime.Gosched()
		}
	}
}
