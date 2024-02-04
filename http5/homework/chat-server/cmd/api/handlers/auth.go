package handlers

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/chatutil"
	"net/http"
)

func (h *Handler) Registration(w http.ResponseWriter, r *http.Request) {
	var newUser chatutil.User
	err := decodeRequestBody(r, &newUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.service.CreateUser(newUser)
	if err != nil {
		http.Error(w, "Create user error", http.StatusInternalServerError)
		return
	}

	err = sendResponse(w, response)
	if err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Authentication(w http.ResponseWriter, r *http.Request) {
	var existingUser chatutil.User
	err := decodeRequestBody(r, &existingUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.service.GetUser(existingUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = sendResponse(w, response)
	if err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
}
