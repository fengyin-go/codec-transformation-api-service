package handler

import (
	"net/http"

	"codec/pkg/httpx"
)

func (s *Server) registerCodecRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/codec/encode", s.encode)
	mux.HandleFunc("POST /api/codec/decode", s.decode)
	mux.HandleFunc("POST /api/codec/batch", s.batchConvert)
	mux.HandleFunc("POST /api/codec/compare", s.compareHash)
	mux.HandleFunc("GET /api/codec/algorithms", s.listSupportedAlgorithms)
	mux.HandleFunc("GET /api/codec/algorithms/{name}/reversible", s.checkReversible)
}

type encodeRequest struct {
	Algorithm string `json:"algorithm"`
	Input     string `json:"input"`
	Key       string `json:"key,omitempty"`
}

func (s *Server) encode(w http.ResponseWriter, r *http.Request) {
	var req encodeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	out, err := s.svc.Encode(req.Algorithm, req.Input, req.Key)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]string{"output": out})
}

type decodeRequest struct {
	Algorithm string `json:"algorithm"`
	Input     string `json:"input"`
	Key       string `json:"key,omitempty"`
}

func (s *Server) decode(w http.ResponseWriter, r *http.Request) {
	var req decodeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	out, err := s.svc.Decode(req.Algorithm, req.Input, req.Key)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]string{"output": out})
}

type batchConvertRequest struct {
	Algorithm string   `json:"algorithm"`
	Mode      string   `json:"mode"`
	Inputs    []string `json:"inputs"`
	Key       string   `json:"key,omitempty"`
}

func (s *Server) batchConvert(w http.ResponseWriter, r *http.Request) {
	var req batchConvertRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	results, err := s.svc.BatchConvertResult(req.Algorithm, req.Mode, req.Inputs, req.Key)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"results": results})
}

type compareHashRequest struct {
	Algorithm string `json:"algorithm"`
	A         string `json:"a"`
	B         string `json:"b"`
	Key       string `json:"key,omitempty"`
}

func (s *Server) compareHash(w http.ResponseWriter, r *http.Request) {
	var req compareHashRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	match, err := s.svc.CompareHash(req.Algorithm, req.A, req.B, req.Key)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]bool{"match": match})
}

func (s *Server) listSupportedAlgorithms(w http.ResponseWriter, r *http.Request) {
	algs := s.svc.SupportedAlgorithms()
	httpx.OK(w, map[string]interface{}{"algorithms": algs})
}

func (s *Server) checkReversible(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	rev := s.svc.IsReversible(name)
	httpx.OK(w, map[string]interface{}{"algorithm": name, "reversible": rev})
}
