package handlers

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/constants"
	"net/http"
	"strconv"
)

func PaginateMessages(messages []entities.Message, limit int, offset int) []entities.Message {
	startIndex := offset
	endIndex := startIndex + limit

	if endIndex > len(messages) {
		endIndex = len(messages)
	}

	return messages[startIndex:endIndex]
}

func GetPaginateParameters(w http.ResponseWriter, r *http.Request) (int, int, error) {
	limitStr := r.URL.Query().Get(constants.LimitQueryParameter)

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, 0, constants.ErrPaginateParameters
	}

	offsetStr := r.URL.Query().Get(constants.OffsetQueryParameter)

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return 0, 0, constants.ErrPaginateParameters
	}

	if limit < 0 || offset < 0 {
		return 0, 0, constants.ErrPaginateParameters
	}

	return limit, offset, nil
}
