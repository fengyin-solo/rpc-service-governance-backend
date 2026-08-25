package model

import (
	"strings"
	"time"
)

const (
	ServiceStatusUp       = "up"
	ServiceStatusDown     = "down"
	ServiceStatusDegraded = "degraded"
)

var serviceTransitions = map[string]map[string]bool{
	ServiceStatusUp:       {ServiceStatusDown: true, ServiceStatusDegraded: true},
	ServiceStatusDown:     {ServiceStatusUp: true},
	ServiceStatusDegraded: {ServiceStatusUp: true},
}

func ServiceCanTransition(from, to string) bool {
	if m, ok := serviceTransitions[from]; ok {
		return m[to]
	}
	return false
}

type Service struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	Owner       string    `json:"owner"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (s *Service) Validate() error {
	s.Name = strings.TrimSpace(s.Name)
	s.Version = strings.TrimSpace(s.Version)
	s.Owner = strings.TrimSpace(s.Owner)
	if s.Name == "" {
		return NewValidationError("name", "服务名称不能为空")
	}
	if s.Version == "" {
		return NewValidationError("version", "版本号不能为空")
	}
	if s.Owner == "" {
		return NewValidationError("owner", "负责人不能为空")
	}
	if s.Status == "" {
		s.Status = ServiceStatusUp
	}
	if s.Status != ServiceStatusUp && s.Status != ServiceStatusDown && s.Status != ServiceStatusDegraded {
		return NewValidationError("status", "服务状态不合法")
	}
	return nil
}

type ServiceFilter struct {
	Status  string
	Owner   string
	Keyword string
}

func (f ServiceFilter) Match(s *Service) bool {
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	if f.Owner != "" && s.Owner != f.Owner {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(s.Name), k) &&
			!strings.Contains(strings.ToLower(s.Description), k) {
			return false
		}
	}
	return true
}
