package handler

import (
	"net/http"

	"codec/internal/model"
	"codec/pkg/httpx"
)

func (s *Server) registerCategoryRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/categories", s.createCategory)
	mux.HandleFunc("GET /api/categories", s.listCategories)
	mux.HandleFunc("GET /api/categories/{id}", s.getCategory)
	mux.HandleFunc("PUT /api/categories/{id}", s.updateCategory)
	mux.HandleFunc("DELETE /api/categories/{id}", s.deleteCategory)
}

type createCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *Server) createCategory(w http.ResponseWriter, r *http.Request) {
	var req createCategoryRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.CreateCategory(model.Category{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

func (s *Server) listCategories(w http.ResponseWriter, r *http.Request) {
	items := s.svc.ListCategories()
	httpx.OK(w, items)
}

func (s *Server) getCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, err := s.svc.GetCategory(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

type updateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *Server) updateCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateCategoryRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.UpdateCategory(id, model.Category{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) deleteCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteCategory(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
