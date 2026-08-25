package store

import (
	"codec/internal/model"
)

func (s *MemoryStore) CreateOperation(o *model.Operation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.operations[o.ID] = o
	return nil
}

func (s *MemoryStore) GetOperation(id string) (*model.Operation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o, ok := s.operations[id]
	if !ok {
		return nil, ErrNotFound
	}
	return o, nil
}

func (s *MemoryStore) ListOperations() []*model.Operation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Operation, 0, len(s.operations))
	for _, o := range s.operations {
		list = append(list, o)
	}
	return list
}

func (s *MemoryStore) DeleteOperation(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.operations[id]; !ok {
		return ErrNotFound
	}
	delete(s.operations, id)
	return nil
}
