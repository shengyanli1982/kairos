package cache

import "sync"

// Segment 结构体
// Segment struct
type Segment struct {
	// lock 是一个互斥锁，用于在多个 goroutine 之间同步对 storage 的访问
	// lock is a mutex lock used to synchronize access to storage across multiple goroutines
	lock sync.Mutex

	// storage 是一个 map，用于存储键值对
	// storage is a map used to store key-value pairs
	storage map[string]any
}

// NewSegment 函数创建并返回一个新的 Segment 实例
// The NewSegment function creates and returns a new Segment instance
func NewSegment() *Segment {
	return &Segment{
		// 初始化互斥锁
		// Initialize the mutex lock
		lock: sync.Mutex{},

		// 创建一个新的 map
		// Create a new map
		storage: make(map[string]any),
	}
}

// Get 方法根据给定的键从 storage 中获取值
// The Get method gets the value from storage based on the given key
func (s *Segment) Get(key string) (any, bool) {
	// 加锁以同步访问
	// Lock to synchronize access
	s.lock.Lock()
	defer s.lock.Unlock()

	// 从 storage 中获取值
	// Get the value from storage
	value, exists := s.storage[key]

	// 返回获取到的值和一个布尔值，该布尔值表示是否找到了该键
	// Return the obtained value and a boolean value, which indicates whether the key was found
	return value, exists
}

// Set 方法将给定的键值对设置到 storage 中
// The Set method sets the given key-value pair to storage
func (s *Segment) Set(key string, value any) {
	// 加锁以同步访问
	// Lock to synchronize access
	s.lock.Lock()
	defer s.lock.Unlock()

	// 将键值对设置到 storage 中
	// Set the key-value pair to storage
	s.storage[key] = value
}

// Delete 方法从 storage 中删除给定的键
// The Delete method deletes the given key from storage
func (s *Segment) Delete(key string) {
	// 加锁以同步访问
	// Lock to synchronize access
	s.lock.Lock()
	defer s.lock.Unlock()

	// 从 storage 中删除键
	// Delete the key from storage
	delete(s.storage, key)
}

// Count 方法返回 storage 中的键值对数量
// The Count method returns the number of key-value pairs in storage
func (s *Segment) Count() int {
	// 加锁以同步访问
	// Lock to synchronize access
	s.lock.Lock()
	defer s.lock.Unlock()

	// 返回 storage 中的键值对数量
	// Return the number of key-value pairs in storage
	return len(s.storage)
}

// Cleanup 方法遍历 storage 中的所有键值对，并对每个值执行给定的函数，然后删除该键值对
// The Cleanup method traverses all key-value pairs in storage, performs the given function on each value, and then deletes the key-value pair
func (s *Segment) Cleanup(fn func(any)) {
	// 加锁以同步访问
	// Lock to synchronize access
	s.lock.Lock()
	defer s.lock.Unlock()

	// 遍历 storage 中的所有键值对
	// Traverse all key-value pairs in storage
	for key, value := range s.storage {
		// 对每个值执行给定的函数
		// Perform the given function on each value
		fn(value)

		// 删除该键值对
		// Delete the key-value pair
		delete(s.storage, key)
	}
}
