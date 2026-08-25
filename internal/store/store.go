// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"rpcgate/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// Service
	CreateService(s *model.Service) error
	GetService(id string) (*model.Service, error)
	ListServices() []*model.Service
	UpdateService(s *model.Service) error
	DeleteService(id string) error

	// Method
	CreateMethod(m *model.Method) error
	GetMethod(id string) (*model.Method, error)
	ListMethods() []*model.Method
	UpdateMethod(m *model.Method) error
	DeleteMethod(id string) error
	DeleteMethodsByServiceID(serviceID string) error

	// Node
	CreateNode(n *model.Node) error
	GetNode(id string) (*model.Node, error)
	ListNodes() []*model.Node
	UpdateNode(n *model.Node) error
	DeleteNode(id string) error
	DeleteNodesByServiceID(serviceID string) error

	// Invocation
	CreateInvocation(i *model.Invocation) error
	GetInvocation(id string) (*model.Invocation, error)
	ListInvocations() []*model.Invocation
	DeleteInvocation(id string) error
	DeleteInvocationsByIDs(ids []string) error

	// TimeoutRule
	CreateTimeoutRule(t *model.TimeoutRule) error
	GetTimeoutRule(id string) (*model.TimeoutRule, error)
	ListTimeoutRules() []*model.TimeoutRule
	UpdateTimeoutRule(t *model.TimeoutRule) error
	DeleteTimeoutRule(id string) error
	DeleteTimeoutRulesByServiceID(serviceID string) error
}
