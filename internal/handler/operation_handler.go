package handler

import (
	"net/http"

	"codec/internal/model"
	"codec/pkg/httpx"
)

func (s *Server) registerOperationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/operations", s.listOperations)
	mux.HandleFunc("GET /api/operations/{id}", s.getOperation)
	mux.HandleFunc("DELETE /api/operations/{id}", s.deleteOperation)
}

func (s *Server) listOperations(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.OperationFilter{
		Algorithm: r.URL.Query().Get("algorithm"),
		Mode:      r.URL.Query().Get("mode"),
	}
	items, total, err := s.svc.ListOperations(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getOperation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	o, err := s.svc.GetOperation(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, o)
}

func (s *Server) deleteOperation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteOperation(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
