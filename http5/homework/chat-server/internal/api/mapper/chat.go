package mapper

import (
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
)

func MakeChatMetadata(username1, username2 string) entities.ChatMetadata {
	if username1 > username2 {
		username1, username2 = username2, username1
	}

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
