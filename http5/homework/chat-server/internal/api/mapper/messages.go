package mapper

import (
	"github.com/ew0s/ewos-to-go-hw/internal/api/resposne"
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
)

func MakeGettingMessagesResponse(messages []entities.Message) resposne.MessagesResponse {
	return resposne.MessagesResponse{
		Messages: messages,
	}
}
