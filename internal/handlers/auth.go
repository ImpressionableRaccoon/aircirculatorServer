package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/storage"
)

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil || len(b) == 0 {
		h.httpJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	var request AuthRequest

	err = json.Unmarshal(b, &request)
	if err != nil {
		h.httpJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	token, err := h.s.SignUp(r.Context(), request.Login, request.Password)
	if err == storage.ErrUserAlreadyExists {
		h.httpJSONError(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		h.httpJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	h.httpJSONStatusOK(w, 201)
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil || len(b) == 0 {
		h.httpJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	var request AuthRequest

	err = json.Unmarshal(b, &request)
	if err != nil {
		h.httpJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	token, err := h.s.SignIn(r.Context(), request.Login, request.Password)
	if err == storage.ErrUnauthorized {
		h.httpJSONError(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		h.httpJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	h.httpJSONStatusOK(w, 200)
}
