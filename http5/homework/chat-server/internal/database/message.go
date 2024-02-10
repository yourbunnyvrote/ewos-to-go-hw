package database

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/constants"
)

func (db *ChatDB) AddPublicMessage(msg entities.Message) error {
	if msg.Content == "" {
		return constants.ErrMsgIsEmpty
	}

	publicChat := db.Get("public chats")

	messages, ok := publicChat.(PublicChatsData)
	if !ok {
		return constants.ErrDataError
	}

	messages = append(messages, msg)

	return nil
}

func (db *ChatDB) GetPublicMessages() ([]entities.Message, error) {
	publicChat := db.Get("public chats")

	messages, ok := publicChat.(PublicChatsData)
	if !ok {
		return nil, constants.ErrDataError
	}

	return messages, nil
}

func (db *ChatDB) AddPrivateMessage(chat entities.Chat, msg entities.Message) error {
	if msg.Content == "" {
		return constants.ErrMsgIsEmpty
	}

	publicChat := db.Get("private chats")

	chats, ok := publicChat.(PrivateChatsData)
	if !ok {
		return constants.ErrDataError
	}

	if chat.User1 < chat.User2 {
		chat.User1, chat.User2 = chat.User2, chat.User1
	}

	chats[chat] = append(chats[chat], msg)

	return nil
}

func (db *ChatDB) GetPrivateMessages(chat entities.Chat) ([]entities.Message, error) {
	privateChat := db.Get("private chats")

	chats, ok := privateChat.(PrivateChatsData)
	if !ok {
		return nil, constants.ErrDataError
	}

	if chat.User1 < chat.User2 {
		chat.User1, chat.User2 = chat.User2, chat.User1
	}

	return chats[chat], nil
}

func (db *ChatDB) GetUsersPrivateMessages(username string) ([]string, error) {
	publicChat := db.Get("private chats")

	chats, ok := publicChat.(PrivateChatsData)
	if !ok {
		return nil, constants.ErrDataError
	}

	result := make([]string, 0)

	for key := range chats {
		if key.User1 == username || key.User2 == username {
			if key.User1 == username {
				result = append(result, key.User2)
			} else {
				result = append(result, key.User1)
			}
		}
	}

	return result, nil
}
