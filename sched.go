package kairos

import (
	"context"
	"sync"
	"time"

	"github.com/shengyanli1982/kairos/internal/cache"
)

type Scheduler struct {
	cfg          *Config
	taskCache    *cache.Cache
	taskCtxCache *cache.Cache
	ctx          context.Context
	cancel       context.CancelFunc
	once         sync.Once
}

func New(conf *Config) *Scheduler {
	conf = isConfigValid(conf)

	s := &Scheduler{
		cfg:          conf,
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

		s.taskCache.Cleanup(func(data any) {
			data.(*Task).Cancel()
		})

		s.taskCtxCache.Cleanup(func(data any) {
			data.(context.CancelFunc)()
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

	s.cfg.callback.OnTaskAdd(id, name, execAt)

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
	taskName := ""

	if task, ok := s.taskCache.Get(id); ok {
		taskName = task.(*Task).GetMetadata().GetName()
		s.taskCache.Delete(id)
		task.(*Task).Wait()
	}

	if cancel, ok := s.taskCtxCache.Get(id); ok {
		cancel.(context.CancelFunc)()
		s.taskCtxCache.Delete(id)
	}

	if taskName != "" {
		s.cfg.callback.OnTaskRemoved(id, taskName)
	}
}
