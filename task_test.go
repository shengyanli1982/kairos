package kairos

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testStandardTaskCallback struct {
	t *testing.T
}

func (tc *testStandardTaskCallback) OnExecuted(id, name string, data any, reason, err error) {
	// handle task callback logic here
	fmt.Printf("Task executed, id: %s, name: %s, data: %v, reason: %v, err: %v\n", id, name, data, reason, err)
	assert.ErrorIs(tc.t, reason, ErrorTaskTimeout, "task should be canceled by timeout")
}

type testEarlyStopTaskCallback struct {
	t *testing.T
}

func (tc *testEarlyStopTaskCallback) OnExecuted(id, name string, data any, reason, err error) {
	// handle task callback logic here
	fmt.Printf("Task executed, id: %s, name: %s, data: %v, reason: %v, err: %v\n", id, name, data, reason, err)
	assert.ErrorIs(tc.t, reason, ErrorTaskEarlyReturn, "task should be canceled by self ctx cancel")
}

type testParentCancelTaskCallback struct {
	t *testing.T
}

func (tc *testParentCancelTaskCallback) OnExecuted(id, name string, data any, reason, err error) {
	// handle task callback logic here
	fmt.Printf("Task executed, id: %s, name: %s, data: %v, reason: %v, err: %v\n", id, name, data, reason, err)
	assert.ErrorIs(tc.t, reason, ErrorTaskCanceled, "task should be canceled by parent ctx cancel")
}

func TestTask_Standard(t *testing.T) {
	t.Run("timeout is less than waiting", func(t *testing.T) {
		parentCtx, parentCancel := context.WithTimeout(context.Background(), time.Millisecond*200)
		defer parentCancel()

		name := "standard task"
		handleFunc := func(done WaitForCtxDone) (any, error) {
			// handle task logic here
			return "lee", nil
		}

		task := NewTask(parentCtx, name, handleFunc, &testStandardTaskCallback{t: t})
		defer task.EarlyReturn()

		time.Sleep(time.Millisecond * 500)
	})

	t.Run("timeout is greater than waiting", func(t *testing.T) {
		parentCtx, parentCancel := context.WithTimeout(context.Background(), time.Millisecond*1000)
		defer parentCancel()

		name := "standard task"
		handleFunc := func(done WaitForCtxDone) (any, error) {
			// handle task logic here
			return "lee", nil
		}

		task := NewTask(parentCtx, name, handleFunc, &testEarlyStopTaskCallback{t: t})
		defer task.EarlyReturn()

		// timeout ctx not cancel, task should be executed after waiting. so trigger the early stop by self ctx
		time.Sleep(time.Millisecond * 500)
	})
}

func TestTask_EarlyStop(t *testing.T) {
	t.Run("timeout is less than waiting", func(t *testing.T) {
		parentCtx, parentCancel := context.WithTimeout(context.Background(), time.Millisecond*200)
		defer parentCancel()

		name := "early stop task"
		handleFunc := func(done WaitForCtxDone) (any, error) {
			// handle task logic here
			return "lee", nil
		}

		task := NewTask(parentCtx, name, handleFunc, &testEarlyStopTaskCallback{t: t})
		task.EarlyReturn()

		time.Sleep(time.Millisecond * 500)
	})

	t.Run("timeout is greater than waiting", func(t *testing.T) {
		parentCtx, parentCancel := context.WithTimeout(context.Background(), time.Millisecond*1000)
		defer parentCancel()

		name := "early stop task"
		handleFunc := func(done WaitForCtxDone) (any, error) {
			// handle task logic here
			return "lee", nil
		}

		task := NewTask(parentCtx, name, handleFunc, &testEarlyStopTaskCallback{t: t})
		task.EarlyReturn()

		time.Sleep(time.Millisecond * 500)
	})
}

func TestTask_ParentCancel(t *testing.T) {
	t.Run("timeout is less than waiting", func(t *testing.T) {
		parentCtx, parentCancel := context.WithTimeout(context.Background(), time.Millisecond*200)

		name := "parent cancel task"
		handleFunc := func(done WaitForCtxDone) (any, error) {
			// handle task logic here
			return "lee", nil
		}

		task := NewTask(parentCtx, name, handleFunc, &testParentCancelTaskCallback{t: t})
		defer task.EarlyReturn()

		parentCancel()

		time.Sleep(time.Millisecond * 500)
	})

	t.Run("timeout is greater than waiting", func(t *testing.T) {
		parentCtx, parentCancel := context.WithTimeout(context.Background(), time.Millisecond*1000)

		name := "parent cancel task"
		handleFunc := func(done WaitForCtxDone) (any, error) {
			// handle task logic here
			return "lee", nil
		}

		task := NewTask(parentCtx, name, handleFunc, &testParentCancelTaskCallback{t: t})
		defer task.EarlyReturn()

		parentCancel()

		time.Sleep(time.Millisecond * 500)
	})
}

func TestTask_WaitingForCtxDone(t *testing.T) {
	parentCtx, parentCancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer parentCancel()

	name := "waiting for ctx done task"
	handleFunc := func(done WaitForCtxDone) (any, error) {
		for {
			select {
			case <-done:
				return "lee", nil
			default:
				time.Sleep(time.Millisecond * 50)
			}
		}
	}

	task := NewTask(parentCtx, name, handleFunc, &testEarlyStopTaskCallback{t: t})
	task.EarlyReturn()

	time.Sleep(time.Millisecond * 500)
}

func TestTask_OnFinished(t *testing.T) {
	t.Run("onFinished function is set correctly", func(t *testing.T) {

		parentCtx, parentCancel := context.WithTimeout(context.Background(), time.Millisecond*200)
		defer parentCancel()

		name := "onfinished task"
		handleFunc := func(done WaitForCtxDone) (any, error) {
			// handle task logic here
			return "lee", nil
		}

		finFunc := func(m *TaskMetadata) {
			fmt.Printf("Task finished, id: %s, name: %s\n", m.id, m.name)
		}

		task := NewTask(parentCtx, name, handleFunc, &testStandardTaskCallback{t: t}).OnFinished(finFunc)
		defer task.EarlyReturn()

		time.Sleep(time.Millisecond * 500)

	})
}
