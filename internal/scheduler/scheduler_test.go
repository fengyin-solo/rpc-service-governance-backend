package scheduler_test

import (
	"context"
	"testing"
	"time"

	"rpcgate/internal/scheduler"
	"rpcgate/internal/worker"
)

type contextRunner struct {
	ctx  chan context.Context
	done chan struct{}
}

func newContextRunner() *contextRunner {
	return &contextRunner{ctx: make(chan context.Context, 1), done: make(chan struct{})}
}

func (r *contextRunner) Run(ctx context.Context) { r.ctx <- ctx; <-ctx.Done(); close(r.done) }
func (r *contextRunner) Stop()                   {}
func (r *contextRunner) Done() <-chan struct{}   { return r.done }

func TestCancellationStopsScheduledRetryWork(t *testing.T) {
	loop := worker.NewLoop(5 * time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	go loop.Run(ctx)
	<-loop.Started()
	cancel()
	select {
	case <-loop.Done():
	case <-time.After(80 * time.Millisecond):
		loop.Stop()
		t.Errorf("worker kept running after its request was cancelled")
	}

	runner := newContextRunner()
	s := scheduler.New(runner)
	scheduleCtx, stopSchedule := context.WithCancel(context.Background())
	s.Start(scheduleCtx)
	received := <-runner.ctx
	stopSchedule()
	select {
	case <-received.Done():
	case <-time.After(80 * time.Millisecond):
		t.Errorf("scheduler replaced the request cancellation context")
	}
}
