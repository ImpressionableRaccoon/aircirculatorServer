package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/service"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/storage"
)

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil || len(b) == 0 {
		h.HTTPJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	var request AuthRequest

	err = json.Unmarshal(b, &request)
	if err != nil {
		h.HTTPJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	token, err := h.s.SignUp(r.Context(), request.Login, request.Password)
	if errors.Is(err, service.ErrLoginIsEmpty) ||
		errors.Is(err, service.ErrPasswordIsEmpty) {
		h.HTTPJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if errors.Is(err, storage.ErrUserAlreadyExists) {
		h.HTTPJSONError(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		h.HTTPJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	h.HTTPJSONStatusOK(w, 201)
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil || len(b) == 0 {
		h.HTTPJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	var request AuthRequest

	err = json.Unmarshal(b, &request)
	if err != nil {
		h.HTTPJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	token, err := h.s.SignIn(r.Context(), request.Login, request.Password)
	if errors.Is(err, service.ErrLoginIsEmpty) ||
		errors.Is(err, service.ErrPasswordIsEmpty) {
		h.HTTPJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if errors.Is(err, storage.ErrUnauthorized) {
		h.HTTPJSONError(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		h.HTTPJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	h.HTTPJSONStatusOK(w, 200)
}
