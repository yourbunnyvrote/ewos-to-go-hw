package mapper

import (
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/public_message/response"
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
)

func MakeSendPublicMessageResponse(msg entities.Message) response.SendPublicMessagesResponse {
	return response.SendPublicMessagesResponse{
		Message: msg,
	}
}

func MakeShowPublicMessagesResponse(messages []entities.Message) response.ShowPublicMessageResponse {
	return response.ShowPublicMessageResponse{
		Messages: messages,
	}
}
