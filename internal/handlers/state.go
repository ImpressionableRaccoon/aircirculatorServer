package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/storage"
)

func (h *Handler) PushDeviceStates(w http.ResponseWriter, r *http.Request) {
	device, err := getDevice(r)
	if err != nil {
		log.Printf("unable to parse device: %v", err)
		h.HTTPJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil || len(b) == 0 {
		h.HTTPJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	var states []storage.DeviceState

	err = json.Unmarshal(b, &states)
	if err != nil {
		h.HTTPJSONError(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = h.s.PushDeviceStates(r.Context(), device, states)
	if err != nil {
		h.HTTPJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.HTTPJSONStatusOK(w, http.StatusCreated)
}
