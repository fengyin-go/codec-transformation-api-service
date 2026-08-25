package store

import (
	"codec/internal/model"
)

func (s *MemoryStore) CreateCategory(c *model.Category) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.categories {
		if exist.Name == c.Name {
			return ErrConflict
		}
	}
	s.categories[c.ID] = c
	return nil
}

func (s *MemoryStore) GetCategory(id string) (*model.Category, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.categories[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (s *MemoryStore) GetCategoryByName(name string) (*model.Category, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.categories {
		if c.Name == name {
			return c, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListCategories() []*model.Category {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Category, 0, len(s.categories))
	for _, c := range s.categories {
		list = append(list, c)
	}
	return list
}

func (s *MemoryStore) UpdateCategory(c *model.Category) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.categories[c.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.categories {
		if exist.ID != c.ID && exist.Name == c.Name {
			return ErrConflict
		}
	}
	s.categories[c.ID] = c
	return nil
}

func (s *MemoryStore) DeleteCategory(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.categories[id]; !ok {
		return ErrNotFound
	}
	delete(s.categories, id)
	return nil
}
