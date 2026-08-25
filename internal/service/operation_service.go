package service

import (
	"sort"
	"time"

	"codec/internal/model"
	"codec/pkg/idgen"
)

func (s *Service) CreateOperation(o model.Operation) (*model.Operation, error) {
	if err := o.Validate(); err != nil {
		return nil, err
	}
	o.ID = idgen.Hex()
	if o.CreatedAt.IsZero() {
		o.CreatedAt = time.Now()
	}
	if err := s.store.CreateOperation(&o); err != nil {
		return nil, err
	}
	return &o, nil
}

func (s *Service) GetOperation(id string) (*model.Operation, error) {
	return s.store.GetOperation(id)
}

func (s *Service) ListOperations(filter model.OperationFilter, page, size int) ([]*model.Operation, int, error) {
	all := s.store.ListOperations()
	matched := make([]*model.Operation, 0, len(all))
	for _, o := range all {
		if filter.Match(o) {
			matched = append(matched, o)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Operation{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) DeleteOperation(id string) error {
	return s.store.DeleteOperation(id)
}
