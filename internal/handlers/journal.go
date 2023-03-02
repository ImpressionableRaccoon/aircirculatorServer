package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/storage"
)

func (h *Handler) GetJournal(w http.ResponseWriter, r *http.Request) {
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

	journal, err := h.s.GetJournalWithCheck(r.Context(), user, deviceID)
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

	body, err := json.Marshal(journal)
	if err != nil {
		h.HTTPJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if len(journal) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(body)
		if err != nil {
			log.Printf("write failed: %v", err)
		}
	}
}
