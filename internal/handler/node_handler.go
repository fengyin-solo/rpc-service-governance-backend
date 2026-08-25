package handler

import (
	"net/http"

	"rpcgate/internal/model"
	"rpcgate/pkg/httpx"
)

func (s *Server) registerNodeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/nodes", s.createNode)
	mux.HandleFunc("GET /api/nodes", s.listNodes)
	mux.HandleFunc("GET /api/nodes/{id}", s.getNode)
	mux.HandleFunc("PUT /api/nodes/{id}", s.updateNode)
	mux.HandleFunc("DELETE /api/nodes/{id}", s.deleteNode)
	mux.HandleFunc("POST /api/nodes/batch", s.batchCreateNodes)
	mux.HandleFunc("DELETE /api/nodes/batch/{service_id}", s.batchDeleteNodesByService)
}

type createNodeRequest struct {
	ServiceID string `json:"service_id"`
	Addr      string `json:"addr"`
	Weight    int    `json:"weight"`
	Region    string `json:"region"`
	Healthy   bool   `json:"healthy"`
}

func (s *Server) createNode(w http.ResponseWriter, r *http.Request) {
	var req createNodeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	n, err := s.svc.CreateNode(model.Node{ServiceID: req.ServiceID, Addr: req.Addr, Weight: req.Weight, Region: req.Region, Healthy: req.Healthy})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, n)
}

func (s *Server) listNodes(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.NodeFilter{
		ServiceID: r.URL.Query().Get("service_id"),
		Region:    r.URL.Query().Get("region"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	if v := r.URL.Query().Get("healthy"); v != "" {
		b := v == "true"
		filter.Healthy = &b
	}
	items, total, err := s.svc.ListNodes(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getNode(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	n, err := s.svc.GetNode(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, n)
}

func (s *Server) updateNode(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createNodeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	n, err := s.svc.UpdateNode(id, model.Node{ServiceID: req.ServiceID, Addr: req.Addr, Weight: req.Weight, Region: req.Region, Healthy: req.Healthy})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, n)
}

func (s *Server) deleteNode(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteNode(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type batchCreateNodesRequest struct {
	Nodes []createNodeRequest `json:"nodes"`
}

func (s *Server) batchCreateNodes(w http.ResponseWriter, r *http.Request) {
	var req batchCreateNodesRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	inputs := make([]model.Node, 0, len(req.Nodes))
	for _, n := range req.Nodes {
		inputs = append(inputs, model.Node{ServiceID: n.ServiceID, Addr: n.Addr, Weight: n.Weight, Region: n.Region, Healthy: n.Healthy})
	}
	results, err := s.svc.BatchCreateNodes(inputs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, results)
}

func (s *Server) batchDeleteNodesByService(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PathValue("service_id")
	if err := s.svc.BatchDeleteNodesByServiceID(serviceID); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
