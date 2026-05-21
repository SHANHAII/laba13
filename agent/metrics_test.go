package main

import (
	"sync"
	"testing"
)

func TestCounterInitial(t *testing.T) {
	c := &Counter{}
	if got := c.Get(); got != 0 {
		t.Errorf("initial counter = %d, want 0", got)
	}
}

func TestCounterInc(t *testing.T) {
	c := &Counter{}
	c.Inc()
	c.Inc()
	if got := c.Get(); got != 2 {
		t.Errorf("counter = %d, want 2", got)
	}
}

func TestCounterConcurrent(t *testing.T) {
	c := &Counter{}
	var wg sync.WaitGroup
	const n = 200
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Inc()
		}()
	}
	wg.Wait()
	if got := c.Get(); got != n {
		t.Errorf("counter = %d, want %d", got, n)
	}
}
