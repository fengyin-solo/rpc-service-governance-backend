package store

import (
	"sync"

	"rpcgate/internal/model"
)

type MemoryStore struct {
	mu           sync.RWMutex
	services     map[string]*model.Service
	methods      map[string]*model.Method
	nodes        map[string]*model.Node
	invocations  map[string]*model.Invocation
	timeoutRules map[string]*model.TimeoutRule
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		services:     make(map[string]*model.Service),
		methods:      make(map[string]*model.Method),
		nodes:        make(map[string]*model.Node),
		invocations:  make(map[string]*model.Invocation),
		timeoutRules: make(map[string]*model.TimeoutRule),
	}
}

var _ Store = (*MemoryStore)(nil)
