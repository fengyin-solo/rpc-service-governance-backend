package service

import (
	"sort"
	"time"

	"rpcgate/internal/model"
	"rpcgate/pkg/idgen"
)

func (s *Service) CreateService(input model.Service) (*model.Service, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	sv := &model.Service{
		ID:          idgen.Hex(),
		Name:        input.Name,
		Version:     input.Version,
		Description: input.Description,
		Owner:       input.Owner,
		Status:      input.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.store.CreateService(sv); err != nil {
		return nil, err
	}
	return sv, nil
}

func (s *Service) GetService(id string) (*model.Service, error) {
	return s.store.GetService(id)
}

func (s *Service) ListServices(filter model.ServiceFilter, page, size int) ([]*model.Service, int, error) {
	if size < 1 {
		size = 20
	}
	if max := s.maxPageSize(); size > max {
		size = max
	}
	all := s.store.ListServices()
	matched := make([]*model.Service, 0, len(all))
	for _, sv := range all {
		if filter.Match(sv) {
			matched = append(matched, sv)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].UpdatedAt.After(matched[j].UpdatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Service{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateService(id string, input model.Service) (*model.Service, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	sv, err := s.store.GetService(id)
	if err != nil {
		return nil, err
	}
	if input.Status != sv.Status {
		if !model.ServiceCanTransition(sv.Status, input.Status) {
			return nil, model.NewValidationError("status", "非法的状态流转")
		}
	}
	sv.Name = input.Name
	sv.Version = input.Version
	sv.Description = input.Description
	sv.Owner = input.Owner
	sv.Status = input.Status
	sv.UpdatedAt = time.Now()
	if err := s.store.UpdateService(sv); err != nil {
		return nil, err
	}
	return sv, nil
}

func (s *Service) DeleteService(id string) error {
	if _, err := s.store.GetService(id); err != nil {
		return err
	}
	if err := s.store.DeleteMethodsByServiceID(id); err != nil {
		return err
	}
	if err := s.store.DeleteNodesByServiceID(id); err != nil {
		return err
	}
	if err := s.store.DeleteTimeoutRulesByServiceID(id); err != nil {
		return err
	}
	return s.store.DeleteService(id)
}
