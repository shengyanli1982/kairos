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

var DefaultScheduledTaskHandleFunc ScheduledTaskHandleFunc = func(ctx context.Context) (data any, err error) { return nil, nil }

var taskPool = sync.Pool{New: func() interface{} { return &ScheduledTask{once: new(sync.Once)} }}

type ScheduledTask struct {
	id         string
	name       string
	callback   ScheduledTaskCallback
	handleFunc ScheduledTaskHandleFunc
	parentCtx  context.Context
	ctx        context.Context
	cancel     context.CancelCauseFunc
	once       *sync.Once
	wg         sync.WaitGroup
}

func NewScheduledTask(parentCtx context.Context, name string, handleFunc ScheduledTaskHandleFunc, callback ScheduledTaskCallback) *ScheduledTask {
	if handleFunc == nil {
		handleFunc = DefaultScheduledTaskHandleFunc
	}

	if parentCtx == nil {
		parentCtx = context.Background()
	}

	if callback == nil {
		callback = NewEmptyTaskCallback()
	}

	task := taskPool.Get().(*ScheduledTask)

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

func (st *ScheduledTask) Stop() {
	st.once.Do(func() {

		st.cancel(ErrorTaskEarlyReturn)

		st.wg.Wait()

		taskPool.Put(st)
	})

	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&st.once)), unsafe.Pointer(new(sync.Once)))
}

func (st *ScheduledTask) executor() {

	taskTrigger := time.NewTicker(time.Millisecond * 500)

	defer func() {
		taskTrigger.Stop()
		st.wg.Done()
	}()

	for {
		select {

		case <-st.ctx.Done():

			switch context.Cause(st.ctx) {

			case context.Canceled:
				st.callback.OnExecuted(st.id, st.name, nil, ErrorTaskCanceled, nil)
				st.cancel(context.Canceled)

			case context.DeadlineExceeded:
				result, err := st.handleFunc(st.ctx)
				st.callback.OnExecuted(st.id, st.name, result, ErrorTaskTimeout, err)

			case ErrorTaskEarlyReturn:
				result, err := st.handleFunc(st.ctx)
				st.callback.OnExecuted(st.id, st.name, result, ErrorTaskEarlyReturn, err)
			}

			return

		case <-taskTrigger.C:
			runtime.Gosched()
		}
	}
}
