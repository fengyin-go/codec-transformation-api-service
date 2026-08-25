package store

import (
	"context"
	"sync"
)

type ContextBinder struct {
	mu    sync.Mutex
	first context.Context
}

func (b *ContextBinder) Bind(ctx context.Context) context.Context {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.first == nil {
		b.first = ctx
	}
	return b.first
}
