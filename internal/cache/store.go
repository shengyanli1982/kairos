package cache

import (
	"sync"

	"github.com/cespare/xxhash/v2"
)

const (
	segmentCount   = uint64(1 << 8)
	segmentsOptVal = segmentCount - 1
)

type Cache struct {
	segments []*Segment
}

func NewCache() *Cache {
	segments := make([]*Segment, segmentCount)
	for i := uint64(0); i < segmentCount; i++ {
		segments[i] = NewSegment()
	}
	return &Cache{
		segments: segments,
	}
}

func (c *Cache) Get(key string) (value any, ok bool) {
	return c.segments[xxhash.Sum64String(key)&segmentsOptVal].Get(key)
}

func (c *Cache) Set(key string, value any) {
	c.segments[xxhash.Sum64String(key)&segmentsOptVal].Set(key, value)
}

func (c *Cache) Delete(key string) {
	c.segments[xxhash.Sum64String(key)&segmentsOptVal].Delete(key)
}

func (c *Cache) Count() int {
	count := 0
	for i := uint64(0); i < segmentCount; i++ {
		count += c.segments[i].Count()
	}
	return count
}

func (c *Cache) Cleanup(fn func(value any)) {
	wg := sync.WaitGroup{}
	for i := uint64(0); i < segmentCount; i++ {
		wg.Add(1)
		go func(idx uint64) {
			defer wg.Done()
			c.segments[idx].Cleanup(fn)
		}(i)
	}
	wg.Wait()
}
