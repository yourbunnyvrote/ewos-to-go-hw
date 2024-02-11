package handlers

import (
	"net/http"
	"strconv"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/constants"
)

func PaginateMessages(messages []entities.Message, limit int, offset int) []entities.Message {
	startIndex := offset
	endIndex := startIndex + limit

	if endIndex > len(messages) {
		endIndex = len(messages)
	}

	if startIndex > len(messages) {
		startIndex = len(messages)
	}

	return messages[startIndex:endIndex]
}

func GetPaginateParameters(r *http.Request) (limit int, offset int, err error) {
	limitStr := r.URL.Query().Get(constants.LimitQueryParameter)

	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		return 0, 0, constants.ErrPaginateParameters
	}

	offsetStr := r.URL.Query().Get(constants.OffsetQueryParameter)

	offset, err = strconv.Atoi(offsetStr)
	if err != nil {
		return 0, 0, constants.ErrPaginateParameters
	}

	if limit < 0 || offset < 0 {
		return 0, 0, constants.ErrPaginateParameters
	}

	return limit, offset, nil
}
