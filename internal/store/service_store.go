package store

import (
	"rpcgate/internal/model"
)

func (s *MemoryStore) CreateService(sv *model.Service) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.services {
		if exist.Name == sv.Name && exist.Version == sv.Version {
			return ErrConflict
		}
	}
	s.services[sv.ID] = sv
	return nil
}

func (s *MemoryStore) GetService(id string) (*model.Service, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sv, ok := s.services[id]
	if !ok {
		return nil, ErrNotFound
	}
	return sv, nil
}

func (s *MemoryStore) ListServices() []*model.Service {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Service, 0, len(s.services))
	for _, sv := range s.services {
		list = append(list, sv)
	}
	return list
}

func (s *MemoryStore) UpdateService(sv *model.Service) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.services[sv.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.services {
		if exist.ID != sv.ID && exist.Name == sv.Name && exist.Version == sv.Version {
			return ErrConflict
		}
	}
	s.services[sv.ID] = sv
	return nil
}

func (s *MemoryStore) DeleteService(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.services[id]; !ok {
		return ErrNotFound
	}
	delete(s.services, id)
	return nil
}
