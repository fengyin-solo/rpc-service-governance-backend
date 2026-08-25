package probe

import (
	"context"
	"errors"
)

type Producer struct{}

func (Producer) Stream(ctx context.Context, nodes []string) (<-chan string, <-chan error) {
	results := make(chan string)
	errs := make(chan error)
	go func() {
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
		close(results)
	}()
	return results, errs
}
