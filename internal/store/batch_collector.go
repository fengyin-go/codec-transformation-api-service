package store

import (
	"context"
	"sort"
)

type BatchResult struct {
	Index  int
	Output string
	Error  string
}

func CollectBatch(ctx context.Context, results <-chan BatchResult, failures <-chan BatchResult) ([]BatchResult, error) {
	collected := make([]BatchResult, 0)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case item, ok := <-results:
			if !ok {
				sort.Slice(collected, func(i, j int) bool { return collected[i].Index < collected[j].Index })
				return collected, nil
			}
			collected = append(collected, item)
		}
	}
}
