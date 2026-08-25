package collector_test

import (
	"context"
	"testing"
	"time"

	"rpcgate/internal/collector"
	"rpcgate/internal/probe"
)

func TestCollectReturnsWhenProbeInputIsInvalid(t *testing.T) {
	producer := probe.Producer{}
	results, errs := producer.Stream(context.Background(), []string{""})
	if err := <-errs; err == nil {
		t.Errorf("invalid probe input returned nil producer error")
	}
	select {
	case _, ok := <-results:
		if ok {
			t.Errorf("invalid probe producer emitted a result")
		}
	case <-time.After(80 * time.Millisecond):
		t.Errorf("invalid probe producer left result channel open")
	}

	c := collector.New(producer)
	done := make(chan error, 1)
	go func() {
		_, err := c.Collect(context.Background(), []string{"node-a", ""})
		done <- err
	}()
	select {
	case err := <-done:
		if err == nil {
			t.Errorf("collector returned nil error for invalid node")
		}
	case <-time.After(120 * time.Millisecond):
		t.Errorf("collector did not return after invalid node")
	}
}
