package handlers

import (
	"encoding/json"
	chatutil "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server"
	"net/http"
	"strconv"
)

func (h *Handler) SendPublicMessage(w http.ResponseWriter, r *http.Request) {
	var msg chatutil.Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg.Username = r.Context().Value("username").(string)

	err = h.service.SendMessage(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(msg)
	if err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetPublicMessage(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	perPage := 3
	startIndex := (page - 1) * perPage
	endIndex := page * perPage

	messages, err := h.service.GetMessage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if startIndex >= len(messages) {
		http.Error(w, "No messages on this page", http.StatusNotFound)
		return
	}

	if endIndex > len(messages) {
		endIndex = len(messages)
	}

	pageMessages := messages[startIndex:endIndex]

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pageMessages)
}

func (h *Handler) GetUsersWithMessages(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) SendPrivateMessage(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {

}
