package service

import (
	"sort"
	"time"

	"rpcgate/internal/model"
	"rpcgate/pkg/idgen"
)

func (s *Service) CreateNode(input model.Node) (*model.Node, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetService(input.ServiceID); err != nil {
		return nil, model.NewValidationError("service_id", "所属服务不存在")
	}
	n := &model.Node{
		ID:              idgen.Hex(),
		ServiceID:       input.ServiceID,
		Addr:            input.Addr,
		Weight:          input.Weight,
		Region:          input.Region,
		Healthy:         input.Healthy,
		LastHeartbeatAt: time.Now(),
		CreatedAt:       time.Now(),
	}
	if err := s.store.CreateNode(n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) GetNode(id string) (*model.Node, error) {
	return s.store.GetNode(id)
}

func (s *Service) ListNodes(filter model.NodeFilter, page, size int) ([]*model.Node, int, error) {
	if size < 1 {
		size = 20
	}
	if max := s.maxPageSize(); size > max {
		size = max
	}
	all := s.store.ListNodes()
	matched := make([]*model.Node, 0, len(all))
	for _, n := range all {
		if filter.Match(n) {
			matched = append(matched, n)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Node{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateNode(id string, input model.Node) (*model.Node, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	n, err := s.store.GetNode(id)
	if err != nil {
		return nil, err
	}
	if input.ServiceID != n.ServiceID {
		if _, err := s.store.GetService(input.ServiceID); err != nil {
			return nil, model.NewValidationError("service_id", "所属服务不存在")
		}
	}
	n.ServiceID = input.ServiceID
	n.Addr = input.Addr
	n.Weight = input.Weight
	n.Region = input.Region
	n.Healthy = input.Healthy
	n.LastHeartbeatAt = time.Now()
	if err := s.store.UpdateNode(n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) DeleteNode(id string) error {
	return s.store.DeleteNode(id)
}

func (s *Service) BatchCreateNodes(inputs []model.Node) ([]*model.Node, error) {
	results := make([]*model.Node, 0, len(inputs))
	for _, input := range inputs {
		n, err := s.CreateNode(input)
		if err != nil {
			return results, err
		}
		results = append(results, n)
	}
	return results, nil
}

func (s *Service) BatchDeleteNodesByServiceID(serviceID string) error {
	if _, err := s.store.GetService(serviceID); err != nil {
		return err
	}
	return s.store.DeleteNodesByServiceID(serviceID)
}
