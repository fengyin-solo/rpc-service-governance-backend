package handler

import (
	"net/http"

	"rpcgate/internal/model"
	"rpcgate/pkg/httpx"
)

func (s *Server) registerMethodRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/methods", s.createMethod)
	mux.HandleFunc("GET /api/methods", s.listMethods)
	mux.HandleFunc("GET /api/methods/{id}", s.getMethod)
	mux.HandleFunc("PUT /api/methods/{id}", s.updateMethod)
	mux.HandleFunc("DELETE /api/methods/{id}", s.deleteMethod)
}

type createMethodRequest struct {
	ServiceID    string `json:"service_id"`
	Name         string `json:"name"`
	ParamSchema  string `json:"param_schema"`
	ReturnSchema string `json:"return_schema"`
	TimeoutMs    int    `json:"timeout_ms"`
	RetryCount   int    `json:"retry_count"`
}

func (s *Server) createMethod(w http.ResponseWriter, r *http.Request) {
	var req createMethodRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	m, err := s.svc.CreateMethod(model.Method{ServiceID: req.ServiceID, Name: req.Name, ParamSchema: req.ParamSchema, ReturnSchema: req.ReturnSchema, TimeoutMs: req.TimeoutMs, RetryCount: req.RetryCount})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, m)
}

func (s *Server) listMethods(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.MethodFilter{
		ServiceID: r.URL.Query().Get("service_id"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListMethods(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getMethod(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	m, err := s.svc.GetMethod(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

func (s *Server) updateMethod(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createMethodRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	m, err := s.svc.UpdateMethod(id, model.Method{ServiceID: req.ServiceID, Name: req.Name, ParamSchema: req.ParamSchema, ReturnSchema: req.ReturnSchema, TimeoutMs: req.TimeoutMs, RetryCount: req.RetryCount})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

func (s *Server) deleteMethod(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteMethod(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
