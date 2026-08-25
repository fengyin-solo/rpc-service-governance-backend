package session

import (
	"errors"
	"sync"
)

var ErrExhausted = errors.New("session pool exhausted")

type Manager struct {
	mu         sync.Mutex
	limit      int
	active     int
	committed  int
	rolledBack int
}

type Lease struct {
	manager *Manager
	once    sync.Once
}

func NewManager(limit int) *Manager { return &Manager{limit: limit} }

func (m *Manager) Begin() (*Lease, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.active >= m.limit {
		return nil, ErrExhausted
	}
	m.active++
	return &Lease{manager: m}, nil
}

func (l *Lease) Finish(_ error) error {
	l.once.Do(func() {
		l.manager.mu.Lock()
		defer l.manager.mu.Unlock()
		l.manager.active--
		l.manager.committed++
	})
	return nil
}

func (m *Manager) Counts() (active, committed, rolledBack int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.active, m.committed, m.rolledBack
}
