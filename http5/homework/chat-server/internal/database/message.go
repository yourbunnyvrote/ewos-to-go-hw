package database

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"
)

func (db *ChatDB) AddPublicMessage(msg entities.Message) error {
	if msg.Content == "" {
		return ErrorMsgIsEmpty
	}

	publicChat := db.Get("public chats")

	messages, ok := publicChat.(PublicChatsData)
	if !ok {
		return ErrorDataError
	}

	db.mu.Lock()
	messages = append(messages, msg)
	db.mu.Unlock()

	db.Insert("public chats", messages)

	return nil
}

func (db *ChatDB) GetPublicMessages() ([]entities.Message, error) {
	publicChat := db.Get("public chats")

	messages, ok := publicChat.(PublicChatsData)
	if !ok {
		return nil, ErrorDataError
	}

	return messages, nil
}

func (db *ChatDB) AddPrivateMessage(chat entities.Chat, msg entities.Message) error {
	if msg.Content == "" {
		return ErrorMsgIsEmpty
	}

	publicChat := db.Get("private chats")

	chats, ok := publicChat.(PrivateChatsData)
	if !ok {
		return ErrorDataError
	}

	if chat.User1 < chat.User2 {
		chat.User1, chat.User2 = chat.User2, chat.User1
	}

	db.mu.Lock()
	chats[chat] = append(chats[chat], msg)
	db.mu.Unlock()

	return nil
}

func (db *ChatDB) GetPrivateMessages(chat entities.Chat) ([]entities.Message, error) {
	privateChat := db.Get("private chats")

	chats, ok := privateChat.(PrivateChatsData)
	if !ok {
		return nil, ErrorDataError
	}

	if chat.User1 < chat.User2 {
		chat.User1, chat.User2 = chat.User2, chat.User1
	}

	db.mu.RLock()
	messages := chats[chat]
	db.mu.RUnlock()

	return messages, nil
}

func (db *ChatDB) GetUsersPrivateMessages(username string) ([]string, error) {
	publicChat := db.Get("private chats")

	chats, ok := publicChat.(PrivateChatsData)
	if !ok {
		return nil, ErrorDataError
	}

	listUsers := make([]string, 0)

	for key := range chats {
		if key.User1 == username || key.User2 == username {
			db.mu.Lock()
			if key.User1 == username {
				listUsers = append(listUsers, key.User2)
			} else {
				listUsers = append(listUsers, key.User1)
			}
			db.mu.Unlock()
		}
	}

	return listUsers, nil
}
