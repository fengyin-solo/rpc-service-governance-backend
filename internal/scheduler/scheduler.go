package scheduler

import (
	"context"
	"time"
)

type Runner interface {
	Run(context.Context)
	Stop()
	Done() <-chan struct{}
}

type Scheduler struct{ runner Runner }

func New(runner Runner) *Scheduler { return &Scheduler{runner: runner} }

func (s *Scheduler) Start(ctx context.Context) { go s.runner.Run(ctx) }

func (s *Scheduler) Shutdown(timeout time.Duration) bool {
	s.runner.Stop()
	select {
	case <-s.runner.Done():
		return true
	case <-time.After(timeout):
		return false
	}
}
