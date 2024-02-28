package message

import "github.com/ew0s/ewos-to-go-hw/internal/domain/entities"

func PaginateMessages(messages []entities.Message, params entities.PaginateParam) []entities.Message {
	startIndex := params.Offset
	endIndex := startIndex + params.Limit

	if endIndex > len(messages) {
		endIndex = len(messages)
	}

	if startIndex > len(messages) {
		startIndex = len(messages)
	}

	return messages[startIndex:endIndex]
}
