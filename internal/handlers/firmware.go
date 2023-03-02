package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/service"
)

func (h *Handler) GetFirmware(w http.ResponseWriter, r *http.Request) {
	md5Hash := r.Header.Get("x-ESP32-sketch-md5")
	sha256Hash := r.Header.Get("x-ESP32-sketch-sha256")

	firmware, err := h.s.GetFirmware(md5Hash, sha256Hash)
	if errors.Is(err, service.ErrFirmwareNoUpdates) {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	if err != nil {
		h.HTTPJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\"firmware.bin\"")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(firmware)
	if err != nil {
		log.Printf("write failed: %v", err)
	}
}
