package model

import (
	"strings"
	"time"
)

const (
	InvocationStatusSuccess = "success"
	InvocationStatusTimeout = "timeout"
	InvocationStatusError   = "error"
)

type Invocation struct {
	ID         string    `json:"id"`
	ServiceID  string    `json:"service_id"`
	MethodID   string    `json:"method_id"`
	NodeID     string    `json:"node_id"`
	RequestID  string    `json:"request_id"`
	DurationMs int       `json:"duration_ms"`
	Status     string    `json:"status"`
	ErrorMsg   string    `json:"error_msg"`
	CalledAt   time.Time `json:"called_at"`
}

func (i *Invocation) Validate() error {
	i.ServiceID = strings.TrimSpace(i.ServiceID)
	i.MethodID = strings.TrimSpace(i.MethodID)
	i.NodeID = strings.TrimSpace(i.NodeID)
	i.RequestID = strings.TrimSpace(i.RequestID)
	if i.ServiceID == "" {
		return NewValidationError("service_id", "所属服务ID不能为空")
	}
	if i.MethodID == "" {
		return NewValidationError("method_id", "方法ID不能为空")
	}
	if i.NodeID == "" {
		return NewValidationError("node_id", "节点ID不能为空")
	}
	if i.Status == "" {
		return NewValidationError("status", "调用状态不能为空")
	}
	if i.Status != InvocationStatusSuccess && i.Status != InvocationStatusTimeout && i.Status != InvocationStatusError {
		return NewValidationError("status", "调用状态不合法")
	}
	return nil
}

type InvocationFilter struct {
	ServiceID string
	MethodID  string
	Status    string
	NodeID    string
}

func (f InvocationFilter) Match(i *Invocation) bool {
	if f.ServiceID != "" && i.ServiceID != f.ServiceID {
		return false
	}
	if f.MethodID != "" && i.MethodID != f.MethodID {
		return false
	}
	if f.Status != "" && i.Status != f.Status {
		return false
	}
	if f.NodeID != "" && i.NodeID != f.NodeID {
		return false
	}
	return true
}
