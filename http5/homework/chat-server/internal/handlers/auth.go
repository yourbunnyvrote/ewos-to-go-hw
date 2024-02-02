package handlers

import (
	"encoding/json"
	chatutil "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server"
	"net/http"
)

func (h *Handler) Registration(w http.ResponseWriter, r *http.Request) {
	var newUser chatutil.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.service.CreateUser(newUser)
	if err != nil {
		http.Error(w, "Create user error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Authentication(w http.ResponseWriter, r *http.Request) {

}
