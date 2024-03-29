package cache

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSegment_Set(t *testing.T) {
	segment := NewSegment()

	// Test case 1: Set a key-value pair
	segment.Set("key1", "value1")
	v, _ := segment.Get("key1")
	assert.Equal(t, "value1", v)

	// Test case 2: Set a key-value pair with concurrent access
	var wg sync.WaitGroup
	numRoutines := 10
	lastV := int64(0)
	wg.Add(numRoutines)
	for i := 0; i < numRoutines; i++ {
		v := i
		go func() {
			defer wg.Done()
			segment.Set("key2", v)
			atomic.StoreInt64(&lastV, int64(v))
		}()
	}
	wg.Wait()

	v, _ = segment.Get("key2")
	assert.Equal(t, lastV, int64(v.(int)))
}

func TestSegment_Delete(t *testing.T) {
	segment := NewSegment()

	// Test case 1: Delete an existing key
	segment.Set("key1", "value1")
	segment.Delete("key1")
	v, _ := segment.Get("key1")
	assert.Nil(t, v)

	// Test case 2: Delete a non-existing key
	segment.Delete("key2")
	v, _ = segment.Get("key2")
	assert.Nil(t, v)
}

func TestSegment_Count(t *testing.T) {
	segment := NewSegment()

	// Test case 1: Count with an empty segment
	assert.Equal(t, 0, segment.Count())

	// Test case 2: Count with non-empty segment
	segment.Set("key1", "value1")
	segment.Set("key2", "value2")
	assert.Equal(t, 2, segment.Count())
}
