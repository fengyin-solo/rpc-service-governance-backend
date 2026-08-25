package model

import (
	"strings"
	"time"
)

type Node struct {
	ID              string    `json:"id"`
	ServiceID       string    `json:"service_id"`
	Addr            string    `json:"addr"`
	Weight          int       `json:"weight"`
	Region          string    `json:"region"`
	Healthy         bool      `json:"healthy"`
	LastHeartbeatAt time.Time `json:"last_heartbeat_at"`
	CreatedAt       time.Time `json:"created_at"`
}

func (n *Node) Validate() error {
	n.ServiceID = strings.TrimSpace(n.ServiceID)
	n.Addr = strings.TrimSpace(n.Addr)
	n.Region = strings.TrimSpace(n.Region)
	if n.ServiceID == "" {
		return NewValidationError("service_id", "所属服务ID不能为空")
	}
	if n.Addr == "" {
		return NewValidationError("addr", "节点地址不能为空")
	}
	if n.Weight < 0 {
		return NewValidationError("weight", "权重不能为负数")
	}
	return nil
}

type NodeFilter struct {
	ServiceID string
	Region    string
	Healthy   *bool
	Keyword   string
}

func (f NodeFilter) Match(n *Node) bool {
	if f.ServiceID != "" && n.ServiceID != f.ServiceID {
		return false
	}
	if f.Region != "" && n.Region != f.Region {
		return false
	}
	if f.Healthy != nil && n.Healthy != *f.Healthy {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(n.Addr), k) {
			return false
		}
	}
	return true
}
