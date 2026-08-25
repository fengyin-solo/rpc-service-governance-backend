package handler

import (
	"net/http"

	"rpcgate/internal/model"
	"rpcgate/pkg/httpx"
)

func (s *Server) registerInvocationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/invocations", s.createInvocation)
	mux.HandleFunc("GET /api/invocations", s.listInvocations)
	mux.HandleFunc("GET /api/invocations/{id}", s.getInvocation)
	mux.HandleFunc("DELETE /api/invocations/{id}", s.deleteInvocation)
	mux.HandleFunc("POST /api/invocations/batch-delete", s.batchDeleteInvocations)
	mux.HandleFunc("GET /api/invocations/report", s.invocationReport)
}

type createInvocationRequest struct {
	ServiceID  string `json:"service_id"`
	MethodID   string `json:"method_id"`
	NodeID     string `json:"node_id"`
	RequestID  string `json:"request_id"`
	DurationMs int    `json:"duration_ms"`
	Status     string `json:"status"`
	ErrorMsg   string `json:"error_msg"`
}

func (s *Server) createInvocation(w http.ResponseWriter, r *http.Request) {
	var req createInvocationRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	i, err := s.svc.CreateInvocation(model.Invocation{ServiceID: req.ServiceID, MethodID: req.MethodID, NodeID: req.NodeID, RequestID: req.RequestID, DurationMs: req.DurationMs, Status: req.Status, ErrorMsg: req.ErrorMsg})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, i)
}

func (s *Server) listInvocations(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.InvocationFilter{
		ServiceID: r.URL.Query().Get("service_id"),
		MethodID:  r.URL.Query().Get("method_id"),
		Status:    r.URL.Query().Get("status"),
		NodeID:    r.URL.Query().Get("node_id"),
	}
	items, total, err := s.svc.ListInvocations(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getInvocation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	i, err := s.svc.GetInvocation(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, i)
}

func (s *Server) deleteInvocation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteInvocation(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type batchDeleteInvocationsRequest struct {
	IDs []string `json:"ids"`
}

func (s *Server) batchDeleteInvocations(w http.ResponseWriter, r *http.Request) {
	var req batchDeleteInvocationsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if err := s.svc.BatchDeleteInvocations(req.IDs); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) invocationReport(w http.ResponseWriter, r *http.Request) {
	report, err := s.svc.InvocationReportByMethod()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, report)
}
