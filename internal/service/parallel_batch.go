package service

import (
	"context"
	"sync"

	"codec/internal/store"
)

func ParallelBatch(ctx context.Context, inputs []string) ([]store.BatchResult, error) {
	// 带缓冲：goroutine 把结果投递后即可 wg.Done() 收尾，
	// 不会因接收方一时未轮到对应 case 而阻塞，从而保证 wg.Wait 能可靠返回、两条 channel 被关闭。
	// 这样混入空输入时，空输入带自身错误及时返回，正常输入带转换结果照常返回，整批不卡住。
	results := make(chan store.BatchResult, len(inputs))
	failures := make(chan store.BatchResult, len(inputs))
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
