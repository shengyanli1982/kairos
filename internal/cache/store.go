package cache

import "github.com/cespare/xxhash/v2"

var (
	SegmentCount    = uint64(1 << 8)
	SegementsOptVal = SegmentCount - 1
)

type Cache struct {
	segments []*Segment
}

func NewCache() *Cache {
	segments := make([]*Segment, SegmentCount)
	for i := uint64(0); i < SegmentCount; i++ {
		segments[i] = NewSegment()
	}
	return &Cache{
		segments: segments,
	}
}

func (c *Cache) Get(key string) (any, bool) {
	return c.segments[xxhash.Sum64String(key)&SegementsOptVal].Get(key)
}

func (c *Cache) Set(key string, value any) {
	c.segments[xxhash.Sum64String(key)&SegementsOptVal].Set(key, value)
}

func (c *Cache) Delete(key string) {
	c.segments[xxhash.Sum64String(key)&SegementsOptVal].Delete(key)
}

func (c *Cache) Count() int {
	count := 0
	for i := uint64(0); i < SegmentCount; i++ {
		count += c.segments[i].Count()
	}
	return count
}

func (c *Cache) Cleanup(fn func(any)) {
	for i := uint64(0); i < SegmentCount; i++ {
		c.segments[i].Cleanup(fn)
	}
}
