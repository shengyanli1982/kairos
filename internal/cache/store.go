package cache

import (
	"sync"

	"github.com/cespare/xxhash/v2"
)

// 定义常量
// Define constants
const (
	// segmentCount 是 Segment 的数量
	// segmentCount is the number of Segments
	segmentCount = uint64(1 << 8)

	// segmentsOptVal 是 segmentCount 减 1，用于优化
	// segmentsOptVal is segmentCount minus 1, used for optimization
	segmentsOptVal = segmentCount - 1
)

// Cache 结构体
// Cache struct
type Cache struct {
	// segments 是一个 Segment 的切片
	// segments is a slice of Segments
	segments []*Segment
}

// NewCache 函数创建并返回一个新的 Cache 实例
// The NewCache function creates and returns a new Cache instance
func NewCache() *Cache {
	// 创建一个 Segment 的切片
	// Create a slice of Segments
	segments := make([]*Segment, segmentCount)

	// 为每个 Segment 创建一个新的 Segment 实例
	// Create a new Segment instance for each Segment
	for i := uint64(0); i < segmentCount; i++ {
		segments[i] = NewSegment()
	}

	// 返回一个新的 Cache 实例
	// Return a new Cache instance
	return &Cache{
		segments: segments,
	}
}

// Get 方法根据给定的键从 Cache 中获取值
// The Get method gets the value from Cache based on the given key
func (c *Cache) Get(key string) (value any, ok bool) {
	// 使用 xxhash.Sum64String 函数计算键的哈希值，然后与 segmentsOptVal 进行与操作，得到索引
	// Use the xxhash.Sum64String function to calculate the hash value of the key, then perform a bitwise AND operation with segmentsOptVal to get the index
	return c.segments[xxhash.Sum64String(key)&segmentsOptVal].Get(key)
}

// Set 方法将给定的键值对设置到 Cache 中
// The Set method sets the given key-value pair to Cache
func (c *Cache) Set(key string, value any) {
	// 使用 xxhash.Sum64String 函数计算键的哈希值，然后与 segmentsOptVal 进行与操作，得到索引
	// Use the xxhash.Sum64String function to calculate the hash value of the key, then perform a bitwise AND operation with segmentsOptVal to get the index
	c.segments[xxhash.Sum64String(key)&segmentsOptVal].Set(key, value)
}

// Delete 方法从 Cache 中删除给定的键
// The Delete method deletes the given key from Cache
func (c *Cache) Delete(key string) {
	// 使用 xxhash.Sum64String 函数计算键的哈希值，然后与 segmentsOptVal 进行与操作，得到索引
	// Use the xxhash.Sum64String function to calculate the hash value of the key, then perform a bitwise AND operation with segmentsOptVal to get the index
	c.segments[xxhash.Sum64String(key)&segmentsOptVal].Delete(key)
}

// Count 方法返回 Cache 中的键值对数量
// The Count method returns the number of key-value pairs in Cache
func (c *Cache) Count() int {
	// 定义一个变量 count，用于存储键值对的数量
	// Define a variable count to store the number of key-value pairs
	count := 0

	// 遍历所有的 Segment，并将每个 Segment 中的键值对数量加到 count 上
	// Traverse all Segments and add the number of key-value pairs in each Segment to count
	for i := uint64(0); i < segmentCount; i++ {
		count += c.segments[i].Count()
	}

	// 返回 count
	// Return count
	return count
}

// Cleanup 方法遍历 Cache 中的所有键值对，并对每个值执行给定的函数，然后删除该键值对
// The Cleanup method traverses all key-value pairs in Cache, performs the given function on each value, and then deletes the key-value pair
func (c *Cache) Cleanup(fn func(value any)) {
	// 定义一个 WaitGroup
	// Define a WaitGroup
	wg := sync.WaitGroup{}

	// 遍历所有的 Segment
	// Traverse all Segments
	for i := uint64(0); i < segmentCount; i++ {

		// 增加 WaitGroup 的计数
		// Increase the count of WaitGroup
		wg.Add(1)

		// 在一个新的 goroutine 中执行 Cleanup 方法
		// Execute the Cleanup method in a new goroutine
		go func(idx uint64) {
			// 在函数结束时调用 Done 方法，减少 WaitGroup 的计数
			// Call the Done method at the end of the function to decrease the count of WaitGroup
			defer wg.Done()

			// 对每个 Segment 执行 Cleanup 方法
			// Execute the Cleanup method for each Segment
			c.segments[idx].Cleanup(fn)

		}(i)
	}

	// 等待所有的 goroutine 完成
	// Wait for all goroutines to complete
	wg.Wait()
}