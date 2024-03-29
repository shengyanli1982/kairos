package kairos

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
)

type OnTaskFinishedHandleFunc = func(metadata *TaskMetadata)

var (
	ErrorTaskCanceled    = errors.New("task canceled")
	ErrorTaskTimeout     = errors.New("task timeout")
	ErrorTaskEarlyReturn = errors.New("task early return")
)

var DefaultTaskHandleFunc TaskHandleFunc = func(done WaitForCtxDone) (data any, err error) { return nil, nil }

var taskPool = sync.Pool{New: func() interface{} { return &Task{metadata: &TaskMetadata{}} }}

type TaskMetadata struct {
	id         string
	name       string
	handleFunc TaskHandleFunc
}

func (stm *TaskMetadata) GetID() string {
	return stm.id
}

func (stm *TaskMetadata) GetName() string {
	return stm.name
}

func (stm *TaskMetadata) GetHandleFunc() TaskHandleFunc {
	return stm.handleFunc
}

type Task struct {
	metadata  *TaskMetadata
	callback  TaskCallback
	parentCtx context.Context
	ctx       context.Context
	cancel    context.CancelCauseFunc
	once      *sync.Once
	wg        *sync.WaitGroup
	onFinFunc OnTaskFinishedHandleFunc
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

	task.metadata.id = uuid.NewString()
	task.metadata.name = name
	task.metadata.handleFunc = handleFunc

	task.callback = callback
	task.parentCtx = parentCtx
	task.ctx, task.cancel = context.WithCancelCause(parentCtx)
	task.wg = &sync.WaitGroup{}
	task.once = &sync.Once{}

	task.wg.Add(1)
	go task.executor()

	return task
}

func (st *Task) executor() {

	taskTrigger := time.NewTicker(time.Millisecond * 500)

	defer func() {
		taskTrigger.Stop()

		st.wg.Done()

		taskPool.Put(st)

		if st.onFinFunc != nil {
			st.onFinFunc(st.metadata)
		}
	}()

	for {
		select {

		case <-st.ctx.Done():

			reason := context.Cause(st.ctx)

			switch reason {

			case context.Canceled:
				st.callback.OnExecuted(st.metadata.id, st.metadata.name, nil, ErrorTaskCanceled, nil)

			case context.DeadlineExceeded:
				result, err := st.metadata.handleFunc(st.ctx.Done())
				st.callback.OnExecuted(st.metadata.id, st.metadata.name, result, ErrorTaskTimeout, err)

			case ErrorTaskEarlyReturn:
				result, err := st.metadata.handleFunc(st.ctx.Done())
				st.callback.OnExecuted(st.metadata.id, st.metadata.name, result, ErrorTaskEarlyReturn, err)
			}

			st.cancel(reason)

			return

		case <-taskTrigger.C:
			runtime.Gosched()
		}
	}
}

func (st *Task) EarlyReturn() {
	if st.once != nil {
		st.once.Do(func() {
			st.cancel(ErrorTaskEarlyReturn)
		})
	}
}

func (st *Task) Cancel() {
	if st.once != nil {
		st.once.Do(func() {
			st.cancel(context.Canceled)
		})
	}
}

func (st *Task) GetMetadata() *TaskMetadata {
	return st.metadata
}

func (st *Task) Wait() {
	st.wg.Wait()
}

func (st *Task) OnFinished(fn OnTaskFinishedHandleFunc) *Task {
	st.onFinFunc = fn
	return st
}
