package store

import (
	"rpcgate/internal/model"
)

func (s *MemoryStore) CreateMethod(m *model.Method) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.methods {
		if exist.ServiceID == m.ServiceID && exist.Name == m.Name {
			return ErrConflict
		}
	}
	s.methods[m.ID] = m
	return nil
}

func (s *MemoryStore) GetMethod(id string) (*model.Method, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.methods[id]
	if !ok {
		return nil, ErrNotFound
	}
	return m, nil
}

func (s *MemoryStore) ListMethods() []*model.Method {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Method, 0, len(s.methods))
	for _, m := range s.methods {
		list = append(list, m)
	}
	return list
}

func (s *MemoryStore) UpdateMethod(m *model.Method) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.methods[m.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.methods {
		if exist.ID != m.ID && exist.ServiceID == m.ServiceID && exist.Name == m.Name {
			return ErrConflict
		}
	}
	s.methods[m.ID] = m
	return nil
}

func (s *MemoryStore) DeleteMethod(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.methods[id]; !ok {
		return ErrNotFound
	}
	delete(s.methods, id)
	return nil
}

func (s *MemoryStore) DeleteMethodsByServiceID(serviceID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, m := range s.methods {
		if m.ServiceID == serviceID {
			delete(s.methods, id)
		}
	}
	return nil
}
