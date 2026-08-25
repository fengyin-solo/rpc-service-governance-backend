package model

import (
	"strings"
	"time"
)

type TimeoutRule struct {
	ID         string    `json:"id"`
	ServiceID  string    `json:"service_id"`
	MethodID   string    `json:"method_id"`
	MaxRetries int       `json:"max_retries"`
	BackoffMs  int       `json:"backoff_ms"`
	TimeoutMs  int       `json:"timeout_ms"`
	Enabled    bool      `json:"enabled"`
	CreatedAt  time.Time `json:"created_at"`
}

func (t *TimeoutRule) Validate() error {
	t.ServiceID = strings.TrimSpace(t.ServiceID)
	t.MethodID = strings.TrimSpace(t.MethodID)
	if t.ServiceID == "" {
		return NewValidationError("service_id", "所属服务ID不能为空")
	}
	if t.MethodID == "" {
		return NewValidationError("method_id", "方法ID不能为空")
	}
	if t.MaxRetries < 0 {
		return NewValidationError("max_retries", "最大重试次数不能为负数")
	}
	if t.BackoffMs < 0 {
		return NewValidationError("backoff_ms", "退避时间不能为负数")
	}
	if t.TimeoutMs < 0 {
		return NewValidationError("timeout_ms", "超时时间不能为负数")
	}
	return nil
}

type TimeoutRuleFilter struct {
	ServiceID string
	MethodID  string
	Enabled   *bool
}

func (f TimeoutRuleFilter) Match(t *TimeoutRule) bool {
	if f.ServiceID != "" && t.ServiceID != f.ServiceID {
		return false
	}
	if f.MethodID != "" && t.MethodID != f.MethodID {
		return false
	}
	if f.Enabled != nil && t.Enabled != *f.Enabled {
		return false
	}
	return true
}
