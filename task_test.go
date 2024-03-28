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

func TestScheduledTask_Standard(t *testing.T) {
	t.Run("timeout is less than waiting", func(t *testing.T) {
		parentCtx, parentCancel := context.WithTimeout(context.Background(), time.Millisecond*200)
		defer parentCancel()

		name := "standard task"
		handleFunc := func(ctx context.Context) (any, error) {
			// handle task logic here
			return "lee", nil
		}

		task := NewScheduledTask(parentCtx, name, handleFunc, &testStandardTaskCallback{t: t})
		defer task.Stop()

		time.Sleep(time.Millisecond * 500)
	})

	t.Run("timeout is greater than waiting", func(t *testing.T) {
		parentCtx, parentCancel := context.WithTimeout(context.Background(), time.Millisecond*1000)
		defer parentCancel()

		name := "standard task"
		handleFunc := func(ctx context.Context) (any, error) {
			// handle task logic here
			return "lee", nil
		}

		task := NewScheduledTask(parentCtx, name, handleFunc, &testEarlyStopTaskCallback{t: t})
		defer task.Stop()

		// timeout ctx not cancel, task should be executed after waiting. so trigger the early stop by self ctx
		time.Sleep(time.Millisecond * 500)
	})
}

func TestScheduledTask_EarlyStop(t *testing.T) {
	t.Run("timeout is less than waiting", func(t *testing.T) {
		parentCtx, parentCancel := context.WithTimeout(context.Background(), time.Millisecond*200)
		defer parentCancel()

		name := "early stop task"
		handleFunc := func(ctx context.Context) (any, error) {
			// handle task logic here
			return "lee", nil
		}

		task := NewScheduledTask(parentCtx, name, handleFunc, &testEarlyStopTaskCallback{t: t})
		task.Stop()

		time.Sleep(time.Millisecond * 500)
	})

	t.Run("timeout is greater than waiting", func(t *testing.T) {
		parentCtx, parentCancel := context.WithTimeout(context.Background(), time.Millisecond*1000)
		defer parentCancel()

		name := "early stop task"
		handleFunc := func(ctx context.Context) (any, error) {
			// handle task logic here
			return "lee", nil
		}

		task := NewScheduledTask(parentCtx, name, handleFunc, &testEarlyStopTaskCallback{t: t})
		task.Stop()

		time.Sleep(time.Millisecond * 500)
	})
}

func TestScheduledTask_ParentCancel(t *testing.T) {
	t.Run("timeout is less than waiting", func(t *testing.T) {
		parentCtx, parentCancel := context.WithTimeout(context.Background(), time.Millisecond*200)

		name := "parent cancel task"
		handleFunc := func(ctx context.Context) (any, error) {
			// handle task logic here
			return "lee", nil
		}

		task := NewScheduledTask(parentCtx, name, handleFunc, &testParentCancelTaskCallback{t: t})
		defer task.Stop()

		parentCancel()

		time.Sleep(time.Millisecond * 500)
	})

	t.Run("timeout is greater than waiting", func(t *testing.T) {
		parentCtx, parentCancel := context.WithTimeout(context.Background(), time.Millisecond*1000)

		name := "parent cancel task"
		handleFunc := func(ctx context.Context) (any, error) {
			// handle task logic here
			return "lee", nil
		}

		task := NewScheduledTask(parentCtx, name, handleFunc, &testParentCancelTaskCallback{t: t})
		defer task.Stop()

		parentCancel()

		time.Sleep(time.Millisecond * 500)
	})
}
