package mapper

import (
	"github.com/ew0s/ewos-to-go-hw/internal/api/resposne"
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
)

func MakeChatMetadata(username1, username2 string) entities.ChatMetadata {
	return entities.ChatMetadata{
		Username1: username1,
		Username2: username2,
	}
}

func MakeEntityMessage(sender string, content string) entities.Message {
	return entities.Message{
		Username: sender,
		Content:  content,
	}
}

func MakeSendingMessageResponse(msg entities.Message) resposne.SendingMessageResponse {
	return resposne.SendingMessageResponse{
		Username: msg.Username,
		Content:  msg.Content,
	}
}
