package store

import "sync"

type PayloadCache struct {
	mu   sync.Mutex
	data map[string][]byte
}

func NewPayloadCache() *PayloadCache { return &PayloadCache{data: make(map[string][]byte)} }

func (c *PayloadCache) Put(key string, payload []byte) {
	c.mu.Lock()
	c.data[key] = payload
	c.mu.Unlock()
}

func (c *PayloadCache) Get(key string) []byte {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.data[key]
}
