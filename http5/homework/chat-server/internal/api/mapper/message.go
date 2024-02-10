package mapper

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"
)

func MakeChat(user1, user2 string) entities.Chat {
	return entities.Chat{
		User1: user1,
		User2: user2,
	}
}

func MakeMessage(sender string, content string) entities.Message {
	return entities.Message{
		Username: sender,
		Content:  content,
	}
}
