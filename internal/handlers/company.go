package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/go-chi/chi/v5"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/storage"
)

func (h *Handler) GetUserCompanies(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r)
	if err != nil {
		log.Printf("unable to parse user: %v", err)
		h.HTTPJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	companies, err := h.s.GetUserCompanies(r.Context(), user)
	if err != nil {
		h.HTTPJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(companies)
	if err != nil {
		h.HTTPJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if len(companies) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(body)
		if err != nil {
			log.Printf("write failed: %v", err)
		}
	}
}

func (h *Handler) AddCompany(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r)
	if err != nil {
		log.Printf("unable to parse user: %v", err)
		h.HTTPJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil || len(b) == 0 {
		h.HTTPJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	var request storage.Company

	err = json.Unmarshal(b, &request)
	if err != nil {
		h.HTTPJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	company, err := h.s.AddCompany(r.Context(), user, request.Name, request.TimeOffset)
	if errors.Is(err, storage.ErrCompanyAlreadyExists) {
		h.HTTPJSONError(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		h.HTTPJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(company)
	if err != nil {
		h.HTTPJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(body)
	if err != nil {
		log.Printf("write failed: %v", err)
	}
}

func (h *Handler) GetCompany(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r)
	if err != nil {
		log.Printf("unable to parse user: %v", err)
		h.HTTPJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	param := chi.URLParam(r, "id")
	companyID, err := uuid.Parse(param)
	if err != nil {
		h.HTTPJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	company, err := h.s.GetUserCompany(r.Context(), user, companyID)
	if errors.Is(err, storage.ErrCompanyNotFound) {
		h.HTTPJSONError(w, err.Error(), http.StatusNotFound)
		return
	}
	if errors.Is(err, storage.ErrCompanyNoPermissions) {
		h.HTTPJSONError(w, err.Error(), http.StatusForbidden)
		return
	}
	if err != nil {
		h.HTTPJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(company)
	if err != nil {
		h.HTTPJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		log.Printf("write failed: %v", err)
	}
}

func (h *Handler) GetCompanyDevices(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r)
	if err != nil {
		log.Printf("unable to parse user: %v", err)
		h.HTTPJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	param := chi.URLParam(r, "id")
	companyID, err := uuid.Parse(param)
	if err != nil {
		h.HTTPJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	devices, err := h.s.GetCompanyDevices(r.Context(), user, companyID)
	if errors.Is(err, storage.ErrCompanyNotFound) {
		h.HTTPJSONError(w, err.Error(), http.StatusNotFound)
		return
	}
	if errors.Is(err, storage.ErrCompanyNoPermissions) {
		h.HTTPJSONError(w, err.Error(), http.StatusForbidden)
		return
	}
	if err != nil {
		h.HTTPJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(devices)
	if err != nil {
		h.HTTPJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if len(devices) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(body)
		if err != nil {
			log.Printf("write failed: %v", err)
		}
	}
}
