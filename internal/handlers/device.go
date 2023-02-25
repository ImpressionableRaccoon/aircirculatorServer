package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/storage"
	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/utils"
)

func (h *Handler) AddDevice(w http.ResponseWriter, r *http.Request) {
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

	var request storage.Device

	err = json.Unmarshal(b, &request)
	if err != nil {
		h.HTTPJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	device, err := h.s.AddDevice(r.Context(), user, request.Company, request.Name, request.Resource)
	if errors.Is(err, storage.ErrCompanyNotFound) {
		h.HTTPJSONError(w, err.Error(), http.StatusForbidden)
		return
	}
	if errors.Is(err, storage.ErrCompanyNoPermissions) {
		h.HTTPJSONError(w, err.Error(), http.StatusForbidden)
		return
	}
	if errors.Is(err, storage.ErrDeviceAlreadyExists) {
		h.HTTPJSONError(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		h.HTTPJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(device)
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

func (h *Handler) GetDevice(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r)
	if err != nil {
		log.Printf("unable to parse user: %v", err)
		h.HTTPJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	param := chi.URLParam(r, "id")
	deviceID, err := uuid.Parse(param)
	if err != nil || deviceID == uuid.Nil {
		h.HTTPJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	device, err := h.s.GetDevice(r.Context(), user, deviceID)
	if errors.Is(err, storage.ErrDeviceNotFound) {
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

	body, err := json.Marshal(device)
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

func (h *Handler) GetDeviceInfo(w http.ResponseWriter, r *http.Request) {
	device, err := getDevice(r)
	if err != nil {
		log.Printf("unable to parse device: %v", err)
		h.HTTPJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(device)
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

func (h *Handler) GetDeviceSchedule(w http.ResponseWriter, r *http.Request) {
	device, err := getDevice(r)
	if err != nil {
		log.Printf("unable to parse device: %v", err)
		h.HTTPJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	schedules, err := h.s.GetDeviceSchedules(r.Context(), device)
	if errors.Is(err, storage.ErrDeviceNotFound) {
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

	body, err := json.Marshal(schedules)
	if err != nil {
		h.HTTPJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Sha256-Hash", utils.GetSHA256(body))
	if len(schedules) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(body)
		if err != nil {
			log.Printf("write failed: %v", err)
		}
	}
}
