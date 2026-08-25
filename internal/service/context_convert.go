package service

import (
	"context"

	"codec/internal/store"
)

type ContextEngine interface {
	Transform(context.Context, string) (string, error)
}

type contextResult struct {
	output string
	err    error
}

func ConvertWithContext(ctx context.Context, binder *store.ContextBinder, engine ContextEngine, input string) (string, error) {
	active := binder.Bind(ctx)
	if err := active.Err(); err != nil {
		return "", err
	}
	result := make(chan contextResult, 1)
	go func() {
		output, err := engine.Transform(context.Background(), input)
		result <- contextResult{output: output, err: err}
	}()
	select {
	case <-active.Done():
		return "", active.Err()
	case item := <-result:
		return item.output, item.err
	}
}
