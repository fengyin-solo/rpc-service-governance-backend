package service

import (
	"sort"
	"time"

	"rpcgate/internal/model"
	"rpcgate/pkg/idgen"
)

func (s *Service) CreateTimeoutRule(input model.TimeoutRule) (*model.TimeoutRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetService(input.ServiceID); err != nil {
		return nil, model.NewValidationError("service_id", "所属服务不存在")
	}
	if _, err := s.store.GetMethod(input.MethodID); err != nil {
		return nil, model.NewValidationError("method_id", "方法不存在")
	}
	t := &model.TimeoutRule{
		ID:         idgen.Hex(),
		ServiceID:  input.ServiceID,
		MethodID:   input.MethodID,
		MaxRetries: input.MaxRetries,
		BackoffMs:  input.BackoffMs,
		TimeoutMs:  input.TimeoutMs,
		Enabled:    input.Enabled,
		CreatedAt:  time.Now(),
	}
	if err := s.store.CreateTimeoutRule(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) GetTimeoutRule(id string) (*model.TimeoutRule, error) {
	return s.store.GetTimeoutRule(id)
}

func (s *Service) ListTimeoutRules(filter model.TimeoutRuleFilter, page, size int) ([]*model.TimeoutRule, int, error) {
	if size < 1 {
		size = 20
	}
	if max := s.maxPageSize(); size > max {
		size = max
	}
	all := s.store.ListTimeoutRules()
	matched := make([]*model.TimeoutRule, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.TimeoutRule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateTimeoutRule(id string, input model.TimeoutRule) (*model.TimeoutRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	t, err := s.store.GetTimeoutRule(id)
	if err != nil {
		return nil, err
	}
	if input.ServiceID != t.ServiceID {
		if _, err := s.store.GetService(input.ServiceID); err != nil {
			return nil, model.NewValidationError("service_id", "所属服务不存在")
		}
	}
	if input.MethodID != t.MethodID {
		if _, err := s.store.GetMethod(input.MethodID); err != nil {
			return nil, model.NewValidationError("method_id", "方法不存在")
		}
	}
	t.ServiceID = input.ServiceID
	t.MethodID = input.MethodID
	t.MaxRetries = input.MaxRetries
	t.BackoffMs = input.BackoffMs
	t.TimeoutMs = input.TimeoutMs
	t.Enabled = input.Enabled
	if err := s.store.UpdateTimeoutRule(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) DeleteTimeoutRule(id string) error {
	return s.store.DeleteTimeoutRule(id)
}
