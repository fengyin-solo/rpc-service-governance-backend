package collector

import (
	"context"

	"rpcgate/internal/probe"
)

type Collector struct{ producer probe.Producer }

func New(producer probe.Producer) *Collector { return &Collector{producer: producer} }

func (c *Collector) Collect(ctx context.Context, nodes []string) ([]string, error) {
	results, errs := c.producer.Stream(ctx, nodes)
	collected := make([]string, 0, len(nodes))

	// 必须同时消费 results 和 errs：遇到空地址时 producer 只向 errs 发送，
	// 若不排空 errs，producer 会阻塞在发送上，进而永远不 close(results)，
	// 这里的 range 也会永久阻塞，本批次卡死。
	var firstErr error
	for {
		select {
		case result, ok := <-results:
			if !ok {
				// results 已关闭；排空 errs（非阻塞）以拿到错误后结束。
				if firstErr != nil {
					return collected, firstErr
				}
				for e := range errs {
					if firstErr == nil {
						firstErr = e
					}
				}
				return collected, firstErr
			}
			collected = append(collected, result)
		case err, ok := <-errs:
			if !ok {
				// errs 已关闭（producer 正常结束）；继续排空 results。
				for result := range results {
					collected = append(collected, result)
				}
				return collected, firstErr
			}
			if firstErr == nil {
				firstErr = err
			}
		case <-ctx.Done():
			// drain 两个通道，确保 producer 的 goroutine 能退出。
			for range results {
			}
			for range errs {
			}
			if firstErr == nil {
				return collected, ctx.Err()
			}
			return collected, firstErr
		}
	}
}
