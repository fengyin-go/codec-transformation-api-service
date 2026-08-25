package store

import (
	"sync"

	"codec/internal/model"
)

type MemoryStore struct {
	mu         sync.RWMutex
	algorithms map[string]*model.Algorithm
	operations map[string]*model.Operation
	categories map[string]*model.Category
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		algorithms: make(map[string]*model.Algorithm),
		operations: make(map[string]*model.Operation),
		categories: make(map[string]*model.Category),
	}
}

var _ Store = (*MemoryStore)(nil)
