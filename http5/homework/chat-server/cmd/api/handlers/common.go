package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"
)

const (
	LimitQueryParameter  = "limit"
	OffsetQueryParameter = "offset"
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

func getPageParams(r *http.Request) (int, int, error) {
	limitStr := r.URL.Query().Get(LimitQueryParameter)

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, 0, err
	}

	offsetStr := r.URL.Query().Get(OffsetQueryParameter)
	offset, err := strconv.Atoi(offsetStr)

	return limit, offset, err
}

func getPaginationIndexes(limit, offset int) (startIndex, endIndex int) {
	startIndex = offset - 1
	endIndex = offset - 1 + limit

	return startIndex, endIndex
}

func paginateMessages(messages []entities.Message, startIndex, endIndex int) ([]entities.Message, error) {
	if startIndex >= len(messages) {
		return nil, ErrEndOfPages
	}

	if endIndex > len(messages) {
		endIndex = len(messages)
	}

	return messages[startIndex:endIndex], nil
}
