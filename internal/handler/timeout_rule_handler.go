package handler

import (
	"net/http"

	"rpcgate/internal/model"
	"rpcgate/pkg/httpx"
)

func (s *Server) registerTimeoutRuleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/timeout-rules", s.createTimeoutRule)
	mux.HandleFunc("GET /api/timeout-rules", s.listTimeoutRules)
	mux.HandleFunc("GET /api/timeout-rules/{id}", s.getTimeoutRule)
	mux.HandleFunc("PUT /api/timeout-rules/{id}", s.updateTimeoutRule)
	mux.HandleFunc("DELETE /api/timeout-rules/{id}", s.deleteTimeoutRule)
}

type createTimeoutRuleRequest struct {
	ServiceID  string `json:"service_id"`
	MethodID   string `json:"method_id"`
	MaxRetries int    `json:"max_retries"`
	BackoffMs  int    `json:"backoff_ms"`
	TimeoutMs  int    `json:"timeout_ms"`
	Enabled    bool   `json:"enabled"`
}

func (s *Server) createTimeoutRule(w http.ResponseWriter, r *http.Request) {
	var req createTimeoutRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.CreateTimeoutRule(model.TimeoutRule{ServiceID: req.ServiceID, MethodID: req.MethodID, MaxRetries: req.MaxRetries, BackoffMs: req.BackoffMs, TimeoutMs: req.TimeoutMs, Enabled: req.Enabled})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, t)
}

func (s *Server) listTimeoutRules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.TimeoutRuleFilter{
		ServiceID: r.URL.Query().Get("service_id"),
		MethodID:  r.URL.Query().Get("method_id"),
	}
	if v := r.URL.Query().Get("enabled"); v != "" {
		b := v == "true"
		filter.Enabled = &b
	}
	items, total, err := s.svc.ListTimeoutRules(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getTimeoutRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := s.svc.GetTimeoutRule(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) updateTimeoutRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createTimeoutRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.UpdateTimeoutRule(id, model.TimeoutRule{ServiceID: req.ServiceID, MethodID: req.MethodID, MaxRetries: req.MaxRetries, BackoffMs: req.BackoffMs, TimeoutMs: req.TimeoutMs, Enabled: req.Enabled})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) deleteTimeoutRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteTimeoutRule(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
