package service

import (
	"context"
	"sync"

	"codec/internal/store"
)

func ParallelBatch(ctx context.Context, inputs []string) ([]store.BatchResult, error) {
	results := make(chan store.BatchResult)
	failures := make(chan store.BatchResult)
	var wg sync.WaitGroup
	for index, input := range inputs {
		wg.Add(1)
		go func(index int, input string) {
			defer wg.Done()
			if input == "" {
				failures <- store.BatchResult{Index: index, Error: "input is empty"}
				return
			}
			results <- store.BatchResult{Index: index, Output: "encoded:" + input}
		}(index, input)
	}
	go func() {
		wg.Wait()
		close(results)
		close(failures)
	}()
	return store.CollectBatch(ctx, results, failures)
}
