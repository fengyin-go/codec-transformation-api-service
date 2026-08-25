package service

import (
	"sort"
	"strings"
	"time"

	"codec/internal/model"
	"codec/pkg/idgen"
)

func (s *Service) CreateCategory(c model.Category) (*model.Category, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	c.ID = idgen.Hex()
	c.CreatedAt = time.Now()
	if err := s.store.CreateCategory(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Service) GetCategory(id string) (*model.Category, error) {
	return s.store.GetCategory(id)
}

func (s *Service) GetCategoryByName(name string) (*model.Category, error) {
	return s.store.GetCategoryByName(name)
}

func (s *Service) ListCategories() []*model.Category {
	list := s.store.ListCategories()
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt.Before(list[j].CreatedAt)
	})
	return list
}

func (s *Service) UpdateCategory(id string, updates model.Category) (*model.Category, error) {
	c, err := s.store.GetCategory(id)
	if err != nil {
		return nil, err
	}
	if updates.Name != "" {
		c.Name = strings.TrimSpace(updates.Name)
	}
	if updates.Description != "" {
		c.Description = strings.TrimSpace(updates.Description)
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateCategory(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) DeleteCategory(id string) error {
	return s.store.DeleteCategory(id)
}
