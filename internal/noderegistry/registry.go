package noderegistry

import "sync"

type Registry struct {
	mu    sync.RWMutex
	nodes map[string]int
}

func New() *Registry { return &Registry{nodes: make(map[string]int)} }

func (r *Registry) Put(address string, weight int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nodes[address] = weight
}

func (r *Registry) Snapshot() map[string]int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.nodes
}
