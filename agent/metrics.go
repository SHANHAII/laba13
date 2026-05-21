package main

import "sync/atomic"

// Counter — потокобезопасный счётчик обработанных задач
type Counter struct {
	n int64
}

func (c *Counter) Inc() {
	atomic.AddInt64(&c.n, 1)
}

func (c *Counter) Get() int64 {
	return atomic.LoadInt64(&c.n)
}
