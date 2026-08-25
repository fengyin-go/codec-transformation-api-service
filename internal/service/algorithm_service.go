package service

import (
	"sort"
	"strings"
	"time"

	"codec/internal/model"
	"codec/pkg/idgen"
)

func (s *Service) CreateAlgorithm(a model.Algorithm) (*model.Algorithm, error) {
	if err := a.Validate(); err != nil {
		return nil, err
	}
	a.ID = idgen.Hex()
	a.CreatedAt = time.Now()
	if err := s.store.CreateAlgorithm(&a); err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *Service) GetAlgorithm(id string) (*model.Algorithm, error) {
	return s.store.GetAlgorithm(id)
}

func (s *Service) GetAlgorithmByName(name string) (*model.Algorithm, error) {
	return s.store.GetAlgorithmByName(name)
}

func (s *Service) ListAlgorithms(filter model.AlgorithmFilter, page, size int) ([]*model.Algorithm, int, error) {
	all := s.store.ListAlgorithms()
	matched := make([]*model.Algorithm, 0, len(all))
	for _, a := range all {
		if filter.Match(a) {
			matched = append(matched, a)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Algorithm{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateAlgorithm(id string, updates model.Algorithm) (*model.Algorithm, error) {
	a, err := s.store.GetAlgorithm(id)
	if err != nil {
		return nil, err
	}
	if updates.Name != "" {
		a.Name = strings.TrimSpace(updates.Name)
	}
	if updates.Category != "" {
		a.Category = strings.TrimSpace(updates.Category)
	}
	if updates.Description != "" {
		a.Description = strings.TrimSpace(updates.Description)
	}
	if err := a.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateAlgorithm(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) DeleteAlgorithm(id string) error {
	return s.store.DeleteAlgorithm(id)
}
