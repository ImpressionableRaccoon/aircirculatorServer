package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r)
	if err != nil {
		log.Printf("unable to parse user: %v", err)
		h.HttpJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(user)
	if err != nil {
		h.HttpJSONError(w, err.Error(), http.StatusInternalServerError)
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
