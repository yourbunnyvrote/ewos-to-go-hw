package handlers

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/constants"
	"net/http"
	"strconv"
)

func PaginateMessages(messages []entities.Message, limit int, offset int) ([]entities.Message, error) {
	startIndex := offset - 1
	endIndex := offset - 1 + limit

	if startIndex >= len(messages) {
		return nil, constants.ErrEndOfPages
	}

	if endIndex > len(messages) {
		endIndex = len(messages)
	}

	return messages[startIndex:endIndex], nil
}

func GetPaginateParameters(w http.ResponseWriter, r *http.Request) (int, int, error) {
	limitStr := r.URL.Query().Get(constants.LimitQueryParameter)

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 0, 0, err
	}

	offsetStr := r.URL.Query().Get(constants.OffsetQueryParameter)

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || limit <= 0 || offset <= 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 0, 0, err
	}

	return limit, offset, nil
}
