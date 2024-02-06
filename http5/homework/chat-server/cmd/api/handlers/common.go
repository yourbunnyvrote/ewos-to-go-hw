package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/chatutil"
)

const (
	PageQueryParameter = "page"
	MessagesPerPage    = 3
)

var ErrEndOfPages = errors.New("no messages on this page")

func decodeRequestBody(r *http.Request, requestBody interface{}) error {
	return json.NewDecoder(r.Body).Decode(requestBody)
}

func sendResponse(w http.ResponseWriter, status int, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(response)
}

func getPageParams(r *http.Request) (int, error) {
	pageStr := r.URL.Query().Get(PageQueryParameter)
	page, err := strconv.Atoi(pageStr)

	return page, err
}

func getPaginationIndexes(page int) (startIndex, endIndex int) {
	startIndex = (page - 1) * MessagesPerPage
	endIndex = page * MessagesPerPage

	return startIndex, endIndex
}

func paginateMessages(messages []chatutil.Message, startIndex, endIndex int) ([]chatutil.Message, error) {
	if startIndex >= len(messages) {
		return nil, ErrEndOfPages
	}

	if endIndex > len(messages) {
		endIndex = len(messages)
	}

	return messages[startIndex:endIndex], nil
}
