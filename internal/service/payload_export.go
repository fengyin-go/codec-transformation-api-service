package service

import (
	"sync"

	"codec/internal/store"
)

type PayloadExporter struct {
	cache   *store.PayloadCache
	mu      sync.Mutex
	exports []string
}

func NewPayloadExporter(cache *store.PayloadCache) *PayloadExporter {
	return &PayloadExporter{cache: cache}
}

func (e *PayloadExporter) Queue(key string, payload []byte, release <-chan struct{}) <-chan struct{} {
	e.cache.Put(key, payload)
	done := make(chan struct{})
	go func() {
		<-release
		e.mu.Lock()
		e.exports = append(e.exports, string(payload))
		e.mu.Unlock()
		close(done)
	}()
	return done
}

func (e *PayloadExporter) Exports() []string {
	e.mu.Lock()
	defer e.mu.Unlock()
	return append([]string(nil), e.exports...)
}
