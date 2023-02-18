package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ImpressionableRaccoon/aircirculatorServer/configs"
	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/storage"
)

type Handler struct {
	st  *storage.PsqlStorage
	cfg *configs.Config
}

func NewHandler(s *storage.PsqlStorage, cfg *configs.Config) *Handler {
	h := &Handler{
		st:  s,
		cfg: cfg,
	}

	return h
}

func (h *Handler) httpJSONError(w http.ResponseWriter, error string, code int) {
	jsonError, _ := json.Marshal(struct {
		Error string `json:"error"`
	}{error})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	_, err := w.Write(jsonError)
	if err != nil {
		log.Printf("write failed: %v", err)
	}
}
