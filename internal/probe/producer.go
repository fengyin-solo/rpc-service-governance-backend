package probe

import (
	"context"
	"errors"
)

type Producer struct{}

func (Producer) Stream(ctx context.Context, nodes []string) (<-chan string, <-chan error) {
	results := make(chan string)
	// errs 带缓冲：遇到空地址时 producer 只发一条错误即返回，缓冲保证发送不阻塞，
	// producer 一定走到 defer close，从而 results 必然被关闭。否则只要消费端没有同时
	// 排空 errs，producer 就会卡在 errs<- 上，永远不 close(results)，本批次卡死。
	errs := make(chan error, 1)
	go func() {
		// 无论正常结束、遇到空地址还是上下文取消，都要关闭两个通道，
		// 否则消费端会在 range results 上永久阻塞。
		defer close(results)
		defer close(errs)
		for _, node := range nodes {
			if node == "" {
				errs <- errors.New("probe node address is empty")
				return
			}
			select {
			case <-ctx.Done():
				return
			case results <- node:
			}
		}
	}()
	return results, errs
}
