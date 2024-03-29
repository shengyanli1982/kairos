package cache

import "sync"

type Segment struct {
	mutex   sync.Mutex
	storage map[string]any
}

func NewSegment() *Segment {
	return &Segment{
		mutex:   sync.Mutex{},
		storage: make(map[string]any),
	}
}

func (s *Segment) Get(key string) (any, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	value, exists := s.storage[key]
	return value, exists
}

func (s *Segment) Set(key string, value any) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.storage[key] = value
}

func (s *Segment) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.storage, key)
}

func (s *Segment) Count() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return len(s.storage)
}

func (s *Segment) Cleanup(fn func(any)) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for key, value := range s.storage {
		fn(value)
		delete(s.storage, key)
	}
}
