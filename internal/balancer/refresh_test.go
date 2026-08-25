package balancer_test

import (
	"context"
	"testing"

	"rpcgate/internal/balancer"
	"rpcgate/internal/noderegistry"
)

func TestRefreshUsesAnImmutableNodeSnapshot(t *testing.T) {
	registry := noderegistry.New()
	registry.Put("node-a", 1)
	snapshot := registry.Snapshot()
	registry.Put("node-b", 2)
	if len(snapshot) != 1 {
		t.Errorf("registry snapshot changed to %v after later registration", snapshot)
	}

	refreshRegistry := noderegistry.New()
	refreshRegistry.Put("node-a", 1)
	b := balancer.New(refreshRegistry)
	ready := make(chan struct{})
	proceed := make(chan struct{})
	result := make(chan int, 1)
	go func() {
		total, _ := b.Refresh(context.Background(), ready, proceed)
		result <- total
	}()
	<-ready
	refreshRegistry.Put("node-b", 2)
	close(proceed)
	if total := <-result; total != 1 {
		t.Errorf("in-flight refresh used later node update, total weight=%d want 1", total)
	}
}
