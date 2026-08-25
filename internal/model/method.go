package model

import (
	"strings"
	"time"
)

type Method struct {
	ID           string    `json:"id"`
	ServiceID    string    `json:"service_id"`
	Name         string    `json:"name"`
	ParamSchema  string    `json:"param_schema"`
	ReturnSchema string    `json:"return_schema"`
	TimeoutMs    int       `json:"timeout_ms"`
	RetryCount   int       `json:"retry_count"`
	CreatedAt    time.Time `json:"created_at"`
}

func (m *Method) Validate() error {
	m.Name = strings.TrimSpace(m.Name)
	m.ServiceID = strings.TrimSpace(m.ServiceID)
	if m.ServiceID == "" {
		return NewValidationError("service_id", "所属服务ID不能为空")
	}
	if m.Name == "" {
		return NewValidationError("name", "方法名称不能为空")
	}
	if m.TimeoutMs < 0 {
		return NewValidationError("timeout_ms", "超时时间不能为负数")
	}
	if m.RetryCount < 0 {
		return NewValidationError("retry_count", "重试次数不能为负数")
	}
	return nil
}

type MethodFilter struct {
	ServiceID string
	Keyword   string
}

func (f MethodFilter) Match(m *Method) bool {
	if f.ServiceID != "" && m.ServiceID != f.ServiceID {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(m.Name), k) {
			return false
		}
	}
	return true
}
