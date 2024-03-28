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

var DefaultScheduledTaskHandleFunc ScheduledTaskHandleFunc = func(done WaitForCtxDone) (data any, err error) { return nil, nil }

var taskPool = sync.Pool{New: func() interface{} { return &ScheduledTask{once: new(sync.Once), metadata: &ScheduledTaskMetadata{}} }}

type ScheduledTaskMetadata struct {
	id         string
	name       string
	handleFunc ScheduledTaskHandleFunc
}

func (stm *ScheduledTaskMetadata) GetID() string {
	return stm.id
}

func (stm *ScheduledTaskMetadata) GetName() string {
	return stm.name
}

func (stm *ScheduledTaskMetadata) GetHandleFunc() ScheduledTaskHandleFunc {
	return stm.handleFunc
}

func (stm *ScheduledTaskMetadata) Reset() {
	stm.id = ""
	stm.name = ""
	stm.handleFunc = nil
}

type ScheduledTask struct {
	metadata  *ScheduledTaskMetadata
	callback  ScheduledTaskCallback
	parentCtx context.Context
	ctx       context.Context
	cancel    context.CancelCauseFunc
	once      *sync.Once
	wg        sync.WaitGroup
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

	task.metadata.id = uuid.NewString()
	task.metadata.name = name
	task.metadata.handleFunc = handleFunc

	task.callback = callback
	task.parentCtx = parentCtx
	task.ctx, task.cancel = context.WithCancelCause(parentCtx)
	task.wg = sync.WaitGroup{}

	task.wg.Add(1)
	go task.executor()

	return task
}

func (st *ScheduledTask) Stop() {
	st.once.Do(func() {

		st.cancel(ErrorTaskEarlyReturn)

		st.wg.Wait()

		st.metadata.Reset()

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
				st.callback.OnExecuted(st.metadata.id, st.metadata.name, nil, ErrorTaskCanceled, nil)
				st.cancel(context.Canceled)

			case context.DeadlineExceeded:
				result, err := st.metadata.handleFunc(st.ctx.Done())
				st.callback.OnExecuted(st.metadata.id, st.metadata.name, result, ErrorTaskTimeout, err)

			case ErrorTaskEarlyReturn:
				result, err := st.metadata.handleFunc(st.ctx.Done())
				st.callback.OnExecuted(st.metadata.id, st.metadata.name, result, ErrorTaskEarlyReturn, err)
			}

			return

		case <-taskTrigger.C:
			runtime.Gosched()
		}
	}
}

func (st *ScheduledTask) GetMetadata() *ScheduledTaskMetadata {
	return st.metadata
}
