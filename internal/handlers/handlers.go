package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ImpressionableRaccoon/aircirculatorServer/configs"
	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/service"
)

type Handler struct {
	s   *service.Service
	cfg *configs.Config
}

func NewHandler(s *service.Service, cfg *configs.Config) *Handler {
	h := &Handler{
		s:   s,
		cfg: cfg,
	}

	return h
}

func (h *Handler) HttpJSONError(w http.ResponseWriter, error string, code int) {
	body, _ := json.Marshal(struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}{
		Status: "error",
		Error:  error,
	})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	_, err := w.Write(body)
	if err != nil {
		log.Printf("write failed: %v", err)
	}
}

func (h *Handler) HttpJSONStatusOK(w http.ResponseWriter, code int) {
	body, _ := json.Marshal(struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	_, err := w.Write(body)
	if err != nil {
		log.Printf("write failed: %v", err)
	}
}
