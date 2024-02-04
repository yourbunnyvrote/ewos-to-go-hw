package handlers

import (
	"encoding/json"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/chatutil"
	"net/http"
	"strconv"
)

const MessagesPerPage = 3

func decodeRequestBody(r *http.Request, requestBody interface{}) error {
	return json.NewDecoder(r.Body).Decode(requestBody)
}

func sendResponse(w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func getPageParams(r *http.Request) (int, error) {
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	return page, err
}

func getPaginationIndexes(page int) (int, int) {
	startIndex := (page - 1) * MessagesPerPage
	endIndex := page * MessagesPerPage
	return startIndex, endIndex
}

func handleMessages(w http.ResponseWriter, messages []chatutil.Message, startIndex, endIndex int) {
	if startIndex >= len(messages) {
		http.Error(w, "No messages on this page", http.StatusNotFound)
		return
	}

	if endIndex > len(messages) {
		endIndex = len(messages)
	}

	pageMessages := messages[startIndex:endIndex]

	err := sendResponse(w, pageMessages)
	if err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
}
