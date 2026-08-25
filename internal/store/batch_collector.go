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

// CollectBatch 从结果与失败两条 channel 收集每条输入的处理结果，按 Index 升序返回。
// 两条 channel 任一就绪都会被消费，因此空输入自身的失败不会拖住整批：它会被尽早带出，
// 正常输入的转换结果照常返回。两条 channel 都关闭后即完成收集。
func CollectBatch(ctx context.Context, results <-chan BatchResult, failures <-chan BatchResult) ([]BatchResult, error) {
	collected := make([]BatchResult, 0)
	for results != nil || failures != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case item, ok := <-results:
			if !ok {
				results = nil // 该 channel 已关闭，从 select 中摘除
				continue
			}
			collected = append(collected, item)
		case item, ok := <-failures:
			if !ok {
				failures = nil // 该 channel 已关闭，从 select 中摘除
				continue
			}
			collected = append(collected, item)
		}
	}
	sort.Slice(collected, func(i, j int) bool { return collected[i].Index < collected[j].Index })
	return collected, nil
}
