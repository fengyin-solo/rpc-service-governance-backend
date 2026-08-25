package service

import (
	"sort"
	"time"

	"rpcgate/internal/model"
	"rpcgate/pkg/idgen"
)

func (s *Service) CreateInvocation(input model.Invocation) (*model.Invocation, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetService(input.ServiceID); err != nil {
		return nil, model.NewValidationError("service_id", "所属服务不存在")
	}
	if _, err := s.store.GetMethod(input.MethodID); err != nil {
		return nil, model.NewValidationError("method_id", "方法不存在")
	}
	if _, err := s.store.GetNode(input.NodeID); err != nil {
		return nil, model.NewValidationError("node_id", "节点不存在")
	}
	sv, err := s.store.GetService(input.ServiceID)
	if err != nil {
		return nil, err
	}
	if sv.Status != model.ServiceStatusUp {
		return nil, model.NewValidationError("service_id", "服务未处于 up 状态，无法承接调用")
	}
	i := &model.Invocation{
		ID:         idgen.Hex(),
		ServiceID:  input.ServiceID,
		MethodID:   input.MethodID,
		NodeID:     input.NodeID,
		RequestID:  input.RequestID,
		DurationMs: input.DurationMs,
		Status:     input.Status,
		ErrorMsg:   input.ErrorMsg,
		CalledAt:   time.Now(),
	}
	if err := s.store.CreateInvocation(i); err != nil {
		return nil, err
	}
	return i, nil
}

func (s *Service) GetInvocation(id string) (*model.Invocation, error) {
	return s.store.GetInvocation(id)
}

func (s *Service) ListInvocations(filter model.InvocationFilter, page, size int) ([]*model.Invocation, int, error) {
	if size < 1 {
		size = 20
	}
	if max := s.maxPageSize(); size > max {
		size = max
	}
	all := s.store.ListInvocations()
	matched := make([]*model.Invocation, 0, len(all))
	for _, i := range all {
		if filter.Match(i) {
			matched = append(matched, i)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CalledAt.After(matched[j].CalledAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Invocation{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) DeleteInvocation(id string) error {
	return s.store.DeleteInvocation(id)
}

func (s *Service) BatchDeleteInvocations(ids []string) error {
	return s.store.DeleteInvocationsByIDs(ids)
}

type MethodReport struct {
	MethodID      string  `json:"method_id"`
	TotalCount    int     `json:"total_count"`
	SuccessCount  int     `json:"success_count"`
	SuccessRate   float64 `json:"success_rate"`
	AvgDurationMs float64 `json:"avg_duration_ms"`
}

func (s *Service) InvocationReportByMethod() ([]MethodReport, error) {
	all := s.store.ListInvocations()
	type agg struct {
		total   int
		success int
		durSum  int
	}
	m := make(map[string]*agg)
	for _, i := range all {
		if _, ok := m[i.MethodID]; !ok {
			m[i.MethodID] = &agg{}
		}
		a := m[i.MethodID]
		a.total++
		if i.Status == model.InvocationStatusSuccess {
			a.success++
		}
		a.durSum += i.DurationMs
	}
	reports := make([]MethodReport, 0, len(m))
	for methodID, a := range m {
		rate := 0.0
		if a.total > 0 {
			rate = float64(a.success) / float64(a.total) * 100
		}
		avg := 0.0
		if a.total > 0 {
			avg = float64(a.durSum) / float64(a.total)
		}
		reports = append(reports, MethodReport{
			MethodID:      methodID,
			TotalCount:    a.total,
			SuccessCount:  a.success,
			SuccessRate:   rate,
			AvgDurationMs: avg,
		})
	}
	sort.Slice(reports, func(i, j int) bool {
		return reports[i].TotalCount > reports[j].TotalCount
	})
	return reports, nil
}
