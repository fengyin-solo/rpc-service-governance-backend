package dispatch

import "context"

type Caller interface {
	Call(context.Context, string) error
}

type Dispatcher struct{ caller Caller }

func New(caller Caller) *Dispatcher { return &Dispatcher{caller: caller} }

func (d *Dispatcher) Dispatch(_ context.Context, requestID string) error {
	return d.caller.Call(context.Background(), requestID)
}
