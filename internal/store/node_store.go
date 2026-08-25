package store

import (
	"rpcgate/internal/model"
)

func (s *MemoryStore) CreateNode(n *model.Node) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.nodes {
		if exist.ServiceID == n.ServiceID && exist.Addr == n.Addr {
			return ErrConflict
		}
	}
	s.nodes[n.ID] = n
	return nil
}

func (s *MemoryStore) GetNode(id string) (*model.Node, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.nodes[id]
	if !ok {
		return nil, ErrNotFound
	}
	return n, nil
}

func (s *MemoryStore) ListNodes() []*model.Node {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Node, 0, len(s.nodes))
	for _, n := range s.nodes {
		list = append(list, n)
	}
	return list
}

func (s *MemoryStore) UpdateNode(n *model.Node) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.nodes[n.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.nodes {
		if exist.ID != n.ID && exist.ServiceID == n.ServiceID && exist.Addr == n.Addr {
			return ErrConflict
		}
	}
	s.nodes[n.ID] = n
	return nil
}

func (s *MemoryStore) DeleteNode(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.nodes[id]; !ok {
		return ErrNotFound
	}
	delete(s.nodes, id)
	return nil
}

func (s *MemoryStore) DeleteNodesByServiceID(serviceID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, n := range s.nodes {
		if n.ServiceID == serviceID {
			delete(s.nodes, id)
		}
	}
	return nil
}
