package downstream

import (
	"context"
	"sync/atomic"
	"time"
)

type Recorder struct {
	delay     time.Duration
	completed atomic.Int64
}

func NewRecorder(delay time.Duration) *Recorder { return &Recorder{delay: delay} }

func (r *Recorder) Call(_ context.Context, _ string) error {
	time.Sleep(r.delay)
	r.completed.Add(1)
	return nil
}

func (r *Recorder) Completed() int64 { return r.completed.Load() }
