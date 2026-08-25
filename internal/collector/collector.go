package collector

import (
	"context"

	"rpcgate/internal/probe"
)

type Collector struct{ producer probe.Producer }

func New(producer probe.Producer) *Collector { return &Collector{producer: producer} }

func (c *Collector) Collect(ctx context.Context, nodes []string) ([]string, error) {
	results, _ := c.producer.Stream(ctx, nodes)
	collected := make([]string, 0, len(nodes))
	for result := range results {
		collected = append(collected, result)
	}
	return collected, nil
}
