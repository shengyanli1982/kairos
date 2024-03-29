package kairos

import (
	"context"
	"sync"
	"time"

	"github.com/shengyanli1982/kairos/internal/cache"
)

type Scheduler struct {
	config       *Config
	taskCache    *cache.Cache
	taskCtxCache *cache.Cache
	ctx          context.Context
	cancel       context.CancelFunc
	once         sync.Once
}

func New(conf *Config) *Scheduler {
	conf = isConfigValid(conf)

	s := &Scheduler{
		config:       conf,
		taskCache:    cache.NewCache(),
		taskCtxCache: cache.NewCache(),
		once:         sync.Once{},
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())

	return s
}

func (s *Scheduler) Stop() {
	s.once.Do(func() {
		s.cancel()

		s.taskCtxCache.Cleanup(func(data any) {
			cancel := data.(context.CancelFunc)
			cancel()
		})

		s.taskCache.Cleanup(func(data any) {
			t := data.(*Task)
			t.Cancel()
		})
	})
}

func (s *Scheduler) add(ctx context.Context, cancel context.CancelFunc, name string, handleFunc TaskHandleFunc) string {
	task := NewTask(ctx, name, handleFunc)

	id := task.GetMetadata().GetID()

	s.taskCache.Set(id, task)

	s.taskCtxCache.Set(id, cancel)

	return id
}

func (s *Scheduler) SetAt(name string, handleFunc TaskHandleFunc, execAt time.Time) string {
	ctx, cancel := context.WithDeadline(s.ctx, execAt)

	id := s.add(ctx, cancel, name, handleFunc)

	s.config.callback.OnTaskAdd(id, name, execAt)

	return id
}

func (s *Scheduler) Set(name string, handleFunc TaskHandleFunc, delay time.Duration) string {
	return s.SetAt(name, handleFunc, time.Now().Add(delay))
}

func (s *Scheduler) Get(id string) *Task {
	if task, ok := s.taskCache.Get(id); ok {
		return task.(*Task)
	}
	return nil
}

func (s *Scheduler) Delete(id string) {
	taskName := "unknown"
	if task, ok := s.taskCache.Get(id); ok {
		taskName = task.(*Task).GetMetadata().GetName()
		s.taskCache.Delete(id)
		s.taskCtxCache.Delete(id)
	}
	s.config.callback.OnTaskRemoved(id, taskName)
}
