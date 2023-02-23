package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) GetSchedules(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r)
	if err != nil {
		log.Printf("unable to parse user: %v", err)
		h.HTTPJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	param := chi.URLParam(r, "id")
	deviceID, err := uuid.Parse(param)
	if err != nil {
		h.HTTPJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	schedules, err := h.s.GetSchedules(r.Context(), user, deviceID)
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

func (h *Handler) AddSchedules(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r)
	if err != nil {
		log.Printf("unable to parse user: %v", err)
		h.HTTPJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	param := chi.URLParam(r, "id")
	deviceID, err := uuid.Parse(param)
	if err != nil {
		h.HTTPJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil || len(b) == 0 {
		h.HTTPJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	var inputSchedules []storage.Schedule

	err = json.Unmarshal(b, &inputSchedules)
	if err != nil {
		h.HTTPJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	schedules, err := h.s.AddSchedules(r.Context(), user, deviceID, inputSchedules)
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
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(body)
	if err != nil {
		log.Printf("write failed: %v", err)
	}
}

func (h *Handler) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r)
	if err != nil {
		log.Printf("unable to parse user: %v", err)
		h.HTTPJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	param := chi.URLParam(r, "id")
	deviceID, err := uuid.Parse(param)
	if err != nil {
		h.HTTPJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	param = chi.URLParam(r, "schedule_id")
	scheduleID, err := uuid.Parse(param)
	if err != nil {
		h.HTTPJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = h.s.DeleteSchedule(r.Context(), user, deviceID, scheduleID)
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

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
}
