package store

import (
	"rpcgate/internal/model"
)

func (s *MemoryStore) CreateTimeoutRule(t *model.TimeoutRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.timeoutRules {
		if exist.ServiceID == t.ServiceID && exist.MethodID == t.MethodID {
			return ErrConflict
		}
	}
	s.timeoutRules[t.ID] = t
	return nil
}

func (s *MemoryStore) GetTimeoutRule(id string) (*model.TimeoutRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.timeoutRules[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (s *MemoryStore) ListTimeoutRules() []*model.TimeoutRule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.TimeoutRule, 0, len(s.timeoutRules))
	for _, t := range s.timeoutRules {
		list = append(list, t)
	}
	return list
}

func (s *MemoryStore) UpdateTimeoutRule(t *model.TimeoutRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.timeoutRules[t.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.timeoutRules {
		if exist.ID != t.ID && exist.ServiceID == t.ServiceID && exist.MethodID == t.MethodID {
			return ErrConflict
		}
	}
	s.timeoutRules[t.ID] = t
	return nil
}

func (s *MemoryStore) DeleteTimeoutRule(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.timeoutRules[id]; !ok {
		return ErrNotFound
	}
	delete(s.timeoutRules, id)
	return nil
}

func (s *MemoryStore) DeleteTimeoutRulesByServiceID(serviceID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, t := range s.timeoutRules {
		if t.ServiceID == serviceID {
			delete(s.timeoutRules, id)
		}
	}
	return nil
}
