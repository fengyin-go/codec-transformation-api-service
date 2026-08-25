package store

import (
	"codec/internal/model"
)

func (s *MemoryStore) CreateAlgorithm(a *model.Algorithm) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.algorithms {
		if exist.Name == a.Name {
			return ErrConflict
		}
	}
	s.algorithms[a.ID] = a
	return nil
}

func (s *MemoryStore) GetAlgorithm(id string) (*model.Algorithm, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.algorithms[id]
	if !ok {
		return nil, ErrNotFound
	}
	return a, nil
}

func (s *MemoryStore) GetAlgorithmByName(name string) (*model.Algorithm, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, a := range s.algorithms {
		if a.Name == name {
			return a, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListAlgorithms() []*model.Algorithm {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Algorithm, 0, len(s.algorithms))
	for _, a := range s.algorithms {
		list = append(list, a)
	}
	return list
}

func (s *MemoryStore) UpdateAlgorithm(a *model.Algorithm) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.algorithms[a.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.algorithms {
		if exist.ID != a.ID && exist.Name == a.Name {
			return ErrConflict
		}
	}
	s.algorithms[a.ID] = a
	return nil
}

func (s *MemoryStore) DeleteAlgorithm(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.algorithms[id]; !ok {
		return ErrNotFound
	}
	delete(s.algorithms, id)
	return nil
}
