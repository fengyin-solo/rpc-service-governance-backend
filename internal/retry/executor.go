package retry

import (
	"context"

	"rpcgate/internal/transport"
)

type Executor struct {
	adapter     *transport.Adapter
	maxAttempts int
}

func New(adapter *transport.Adapter, maxAttempts int) *Executor {
	return &Executor{adapter: adapter, maxAttempts: maxAttempts}
}

func (e *Executor) Do(ctx context.Context, payload string) error {
	var last error
	for attempt := 0; attempt < e.maxAttempts; attempt++ {
		last = e.adapter.Call(ctx, payload)
		if last == nil {
			return nil
		}
		if err := ctx.Err(); err != nil {
			return err
		}
	}
	return last
}
