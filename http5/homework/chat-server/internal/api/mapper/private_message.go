package mapper

import (
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/private_message/request"
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/private_message/response"
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
)

func MakeShowPrivateMessageRequest(receiver, sender string) request.ShowPrivateMessageRequest {
	return request.ShowPrivateMessageRequest{
		Receiver: receiver,
		Sender:   sender,
	}
}

func MakeMessagesResponse(messages []entities.Message) []response.Message {
	result := make([]response.Message, len(messages))

	for i, msg := range messages {
		result[i] = response.Message{
			Username: msg.Username,
			Content:  msg.Content,
		}
	}

	return result
}

func MakeShowPrivateMessageResponse(messages []response.Message) response.ShowPrivateMessagesResponse {
	return response.ShowPrivateMessagesResponse{
		Messages: messages,
	}
}

func MakeSendPrivateMessageResponse(msg entities.Message) response.SendPrivateMessageResponse {
	return response.SendPrivateMessageResponse{
		Username: msg.Username,
		Content:  msg.Content,
	}
}
