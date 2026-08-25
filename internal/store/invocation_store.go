package store

import (
	"rpcgate/internal/model"
)

func (s *MemoryStore) CreateInvocation(i *model.Invocation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.invocations[i.ID] = i
	return nil
}

func (s *MemoryStore) GetInvocation(id string) (*model.Invocation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	i, ok := s.invocations[id]
	if !ok {
		return nil, ErrNotFound
	}
	return i, nil
}

func (s *MemoryStore) ListInvocations() []*model.Invocation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Invocation, 0, len(s.invocations))
	for _, i := range s.invocations {
		list = append(list, i)
	}
	return list
}

func (s *MemoryStore) DeleteInvocation(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.invocations[id]; !ok {
		return ErrNotFound
	}
	delete(s.invocations, id)
	return nil
}

func (s *MemoryStore) DeleteInvocationsByIDs(ids []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, id := range ids {
		delete(s.invocations, id)
	}
	return nil
}
