package service

import "sort"

type OverviewStats struct {
	ServiceStatusGroup    map[string]int     `json:"service_status_group"`
	InvocationStatusGroup map[string]int     `json:"invocation_status_group"`
	NodeHealthRatio       map[string]float64 `json:"node_health_ratio"`
	TopSlowMethods        []SlowMethodItem   `json:"top_slow_methods"`
}

type SlowMethodItem struct {
	MethodID      string  `json:"method_id"`
	AvgDurationMs float64 `json:"avg_duration_ms"`
}

func (s *Service) OverviewStats() (*OverviewStats, error) {
	result := &OverviewStats{
		ServiceStatusGroup:    make(map[string]int),
		InvocationStatusGroup: make(map[string]int),
		NodeHealthRatio:       make(map[string]float64),
	}

	for _, sv := range s.store.ListServices() {
		result.ServiceStatusGroup[sv.Status]++
	}

	for _, i := range s.store.ListInvocations() {
		result.InvocationStatusGroup[i.Status]++
	}

	var healthyCount, totalNodes int
	for _, n := range s.store.ListNodes() {
		totalNodes++
		if n.Healthy {
			healthyCount++
		}
	}
	if totalNodes > 0 {
		result.NodeHealthRatio["healthy"] = float64(healthyCount) / float64(totalNodes) * 100
		result.NodeHealthRatio["unhealthy"] = float64(totalNodes-healthyCount) / float64(totalNodes) * 100
	} else {
		result.NodeHealthRatio["healthy"] = 0
		result.NodeHealthRatio["unhealthy"] = 0
	}

	type agg struct {
		durSum int
		count  int
	}
	m := make(map[string]*agg)
	for _, i := range s.store.ListInvocations() {
		if _, ok := m[i.MethodID]; !ok {
			m[i.MethodID] = &agg{}
		}
		m[i.MethodID].durSum += i.DurationMs
		m[i.MethodID].count++
	}
	slow := make([]SlowMethodItem, 0, len(m))
	for methodID, a := range m {
		avg := 0.0
		if a.count > 0 {
			avg = float64(a.durSum) / float64(a.count)
		}
		slow = append(slow, SlowMethodItem{MethodID: methodID, AvgDurationMs: avg})
	}
	sort.Slice(slow, func(i, j int) bool {
		return slow[i].AvgDurationMs > slow[j].AvgDurationMs
	})
	if len(slow) > 10 {
		slow = slow[:10]
	}
	result.TopSlowMethods = slow

	return result, nil
}
