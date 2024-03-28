package kairos

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type Timer struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	once   sync.Once
	timer  atomic.Int64
}

func NewTimer() *Timer {
	t := &Timer{
		wg:    sync.WaitGroup{},
		once:  sync.Once{},
		timer: atomic.Int64{},
	}
	t.ctx, t.cancel = context.WithCancel(context.Background())

	t.wg.Add(1)
	go t.updateTimer()

	return t
}

func (t *Timer) Stop() {
	t.once.Do(func() {
		t.cancel()
		t.wg.Wait()
	})
}

func (t *Timer) Now() int64 {
	return t.timer.Load()
}

// updateTimer 是一个定时器，用于更新 WriteAsyncer 的时间戳
// updateTimer is a timer that updates the timestamp of WriteAsyncer
func (t *Timer) updateTimer() {
	// 创建一个每秒触发一次的定时器
	// Create a timer that triggers once per second
	ticker := time.NewTicker(time.Second)

	// 使用 defer 语句确保在函数退出时停止定时器并减少等待组的计数
	// Use a defer statement to ensure that the timer is stopped and the wait group count is decreased when the function exits
	defer func() {
		ticker.Stop() // 停止定时器
		t.wg.Done()   // 减少等待组的计数
	}()

	// 使用无限循环来持续更新时间戳
	// Use an infinite loop to continuously update the timestamp
	for {
		select {
		// 如果收到上下文的 Done 信号，就退出循环
		// If the Done signal of the context is received, exit the loop
		case <-t.ctx.Done():
			return

		// 如果定时器触发，就更新时间戳
		// If the timer triggers, update the timestamp
		case <-ticker.C:
			// 使用当前的 Unix 毫秒时间戳更新 timer
			// Update timer with the current Unix millisecond timestamp
			t.timer.Store(time.Now().UnixMilli())
		}
	}
}
