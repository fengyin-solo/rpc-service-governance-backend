package balancer

import "context"

type SnapshotSource interface {
	Snapshot() map[string]int
}

type Balancer struct{ source SnapshotSource }

func New(source SnapshotSource) *Balancer { return &Balancer{source: source} }

func (b *Balancer) Refresh(ctx context.Context, ready chan<- struct{}, proceed <-chan struct{}) (int, error) {
	snapshot := b.source.Snapshot()
	close(ready)
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case <-proceed:
	}
	total := 0
	for _, weight := range snapshot {
		total += weight
	}
	return total, nil
}
