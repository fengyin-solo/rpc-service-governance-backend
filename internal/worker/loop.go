package worker

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type Loop struct {
	interval time.Duration
	calls    atomic.Int64
	started  chan struct{}
	stop     chan struct{}
	done     chan struct{}
	once     sync.Once
}

func NewLoop(interval time.Duration) *Loop {
	return &Loop{interval: interval, started: make(chan struct{}), stop: make(chan struct{}), done: make(chan struct{})}
}

func (l *Loop) Run(_ context.Context) {
	defer close(l.done)
	ticker := time.NewTicker(l.interval)
	defer ticker.Stop()
	for {
		select {
		case <-l.stop:
			return
		case <-ticker.C:
			if l.calls.Add(1) == 1 {
				close(l.started)
			}
		}
	}
}

func (l *Loop) Stop()                    { l.once.Do(func() { close(l.stop) }) }
func (l *Loop) Done() <-chan struct{}    { return l.done }
func (l *Loop) Started() <-chan struct{} { return l.started }
func (l *Loop) Calls() int64             { return l.calls.Load() }
