package transport

import (
	"context"
	"errors"
	"fmt"
)

type Backend interface {
	Invoke(context.Context, string) error
}

type RejectedError struct{ Reason string }

func (e *RejectedError) Error() string { return "rpc rejected: " + e.Reason }

type TemporaryError struct{ Reason string }

func (e *TemporaryError) Error() string { return "rpc temporary failure: " + e.Reason }

type Adapter struct{ backend Backend }

func NewAdapter(backend Backend) *Adapter { return &Adapter{backend: backend} }

func (a *Adapter) Call(ctx context.Context, payload string) error {
	if err := a.backend.Invoke(ctx, payload); err != nil {
		return fmt.Errorf("backend invocation: %v", err)
	}
	return nil
}

func IsTemporary(err error) bool {
	var target *TemporaryError
	return errors.As(err, &target)
}

func IsRejected(err error) bool {
	var target *RejectedError
	return errors.As(err, &target)
}
