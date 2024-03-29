package kairos

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Define a struct for test scheduler callback
type testSchedCallback struct{}

// OnTaskExecuted is called when a task is executed
func (tc *testSchedCallback) OnTaskExecuted(id, name string, data any, reason, err error) {
	// Print the task execution details
	fmt.Printf("Task executed, id: %s, name: %s, data: %v, reason: %v, err: %v\n", id, name, data, reason, err)
}

// OnTaskAdded is called when a task is added
func (tc *testSchedCallback) OnTaskAdded(id, name string, execAt time.Time) {
	// Print the task addition details
	fmt.Printf("Task added, id: %s, name: %s, execAt: %s\n", id, name, execAt.String())
}

// OnTaskRemoved is called when a task is removed
func (tc *testSchedCallback) OnTaskRemoved(id, name string) {
	// Print the task removal details
	fmt.Printf("Task removed, id: %s, name: %s\n", id, name)
}

// TestScheduler_Set is a test function for the Set method of the Scheduler
func TestScheduler_Set(t *testing.T) {
	// Create a new configuration with the test scheduler callback
	config := NewConfig().WithCallback(&testSchedCallback{})

	// Create a new scheduler with the configuration
	scheduler := New(config)

	// Add 10 tasks to the scheduler
	for i := 0; i < 10; i++ {
		index := i
		taskID := scheduler.Set("test", func(_ WaitForContextDone) (result any, err error) {
			// The task simply returns a string
			return fmt.Sprintf("test task: %d", index), nil
		}, time.Millisecond*200)

		// Get the task from the scheduler and assert it's not nil
		task := scheduler.Get(taskID)
		assert.NotNil(t, task)
	}

	// Sleep for a while to let the tasks execute
	time.Sleep(time.Millisecond * 500)

	// Stop the scheduler
	scheduler.Stop()

	// Assert that all tasks have been executed and removed from the scheduler
	assert.Equal(t, 0, scheduler.GetTaskCount())
}

// TestScheduler_SetAt is a test function for the SetAt method of the Scheduler
func TestScheduler_SetAt(t *testing.T) {
	// Create a new configuration with the test scheduler callback
	config := NewConfig().WithCallback(&testSchedCallback{})

	// Create a new scheduler with the configuration
	scheduler := New(config)

	// Add 10 tasks to the scheduler with a specific execution time
	for i := 0; i < 10; i++ {
		index := i
		taskID := scheduler.SetAt("test", func(_ WaitForContextDone) (result any, err error) {
			// The task simply returns a string
			return fmt.Sprintf("test task: %d", index), nil
		}, time.Now().Add(time.Millisecond*200))

		// Get the task from the scheduler and assert it's not nil
		task := scheduler.Get(taskID)
		assert.NotNil(t, task)
	}

	// Sleep for a while to let the tasks execute
	time.Sleep(time.Millisecond * 500)

	// Stop the scheduler
	scheduler.Stop()
	// Assert that all tasks have been executed and removed from the scheduler
	assert.Equal(t, 0, scheduler.GetTaskCount())
}

// TestScheduler_Delete is a test function for the Delete method of the Scheduler
func TestScheduler_Delete(t *testing.T) {
	// Create a new configuration with the test scheduler callback
	config := NewConfig().WithCallback(&testSchedCallback{})

	// Create a new scheduler with the configuration
	scheduler := New(config)

	// Add 10 tasks to the scheduler and then delete them
	for i := 0; i < 10; i++ {
		index := i
		taskID := scheduler.SetAt("test", func(_ WaitForContextDone) (result any, err error) {
			// The task simply returns a string
			return fmt.Sprintf("test task: %d", index), nil
		}, time.Now().Add(time.Millisecond*200))

		// Get the task from the scheduler and assert it's not nil
		task := scheduler.Get(taskID)
		assert.NotNil(t, task)

		// Delete the task from the scheduler
		scheduler.Delete(taskID)

		// Get the task from the scheduler again and assert it's now nil
		task = scheduler.Get(taskID)
		assert.Nil(t, task)
	}

	// Sleep for a while to let the tasks execute
	time.Sleep(time.Millisecond * 500)

	// Stop the scheduler
	scheduler.Stop()
	// Assert that all tasks have been executed and removed from the scheduler
	assert.Equal(t, 0, scheduler.GetTaskCount())
}

// TestScheduler_Stop is a test function for the Stop method of the Scheduler
func TestScheduler_Stop(t *testing.T) {
	// Create a new configuration with the test scheduler callback
	config := NewConfig().WithCallback(&testSchedCallback{})

	// Create a new scheduler with the configuration
	scheduler := New(config)

	// Add 10 tasks to the scheduler
	for i := 0; i < 10; i++ {
		index := i
		taskID := scheduler.SetAt("test", func(_ WaitForContextDone) (result any, err error) {
			// The task simply returns a string
			return fmt.Sprintf("test task: %d", index), nil
		}, time.Now().Add(time.Millisecond*200))

		// Get the task from the scheduler and assert it's not nil
		task := scheduler.Get(taskID)
		assert.NotNil(t, task)
	}

	// Stop the scheduler
	scheduler.Stop()
	// Assert that all tasks have been stopped and removed from the scheduler
	assert.Equal(t, 0, scheduler.GetTaskCount())
}
