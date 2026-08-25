package service

import (
	"sort"
	"time"

	"rpcgate/internal/model"
	"rpcgate/pkg/idgen"
)

func (s *Service) CreateMethod(input model.Method) (*model.Method, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetService(input.ServiceID); err != nil {
		return nil, model.NewValidationError("service_id", "所属服务不存在")
	}
	m := &model.Method{
		ID:           idgen.Hex(),
		ServiceID:    input.ServiceID,
		Name:         input.Name,
		ParamSchema:  input.ParamSchema,
		ReturnSchema: input.ReturnSchema,
		TimeoutMs:    input.TimeoutMs,
		RetryCount:   input.RetryCount,
		CreatedAt:    time.Now(),
	}
	if err := s.store.CreateMethod(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) GetMethod(id string) (*model.Method, error) {
	return s.store.GetMethod(id)
}

func (s *Service) ListMethods(filter model.MethodFilter, page, size int) ([]*model.Method, int, error) {
	if size < 1 {
		size = 20
	}
	if max := s.maxPageSize(); size > max {
		size = max
	}
	all := s.store.ListMethods()
	matched := make([]*model.Method, 0, len(all))
	for _, m := range all {
		if filter.Match(m) {
			matched = append(matched, m)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Method{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateMethod(id string, input model.Method) (*model.Method, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	m, err := s.store.GetMethod(id)
	if err != nil {
		return nil, err
	}
	if input.ServiceID != m.ServiceID {
		if _, err := s.store.GetService(input.ServiceID); err != nil {
			return nil, model.NewValidationError("service_id", "所属服务不存在")
		}
	}
	m.ServiceID = input.ServiceID
	m.Name = input.Name
	m.ParamSchema = input.ParamSchema
	m.ReturnSchema = input.ReturnSchema
	m.TimeoutMs = input.TimeoutMs
	m.RetryCount = input.RetryCount
	if err := s.store.UpdateMethod(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) DeleteMethod(id string) error {
	return s.store.DeleteMethod(id)
}
