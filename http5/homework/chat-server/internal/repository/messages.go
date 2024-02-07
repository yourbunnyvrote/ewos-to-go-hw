package repository

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/database"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"
)

type ChatsDB struct {
	db InMemoryDB
}

func NewChatsDB(db InMemoryDB) *ChatsDB {
	return &ChatsDB{db: db}
}

func (pc *ChatsDB) SendPublicMessage(msg entities.Message) error {
	publicChat := pc.db.Get("public chats")

	messages, ok := publicChat.(database.PublicChatsData)
	if !ok {
		return ErrDataError
	}

	messages = append(messages, msg)

	pc.db.Insert("public chats", messages)

	return nil
}

func (pc *ChatsDB) SendPrivateMessage(chat database.Chat, msg entities.Message) error {
	publicChat := pc.db.Get("private chats")

	chats, ok := publicChat.(database.PrivateChatsData)
	if !ok {
		return ErrDataError
	}

	chats[chat] = append(chats[chat], msg)

	pc.db.Insert("private chats", chats)

	return nil
}

func (pc *ChatsDB) GetPublicMessages() ([]entities.Message, error) {
	publicChat := pc.db.Get("public chats")

	messages, ok := publicChat.(database.PublicChatsData)
	if !ok {
		return nil, ErrDataError
	}

	return messages, nil
}

func (pc *ChatsDB) GetPrivateMessages(chat database.Chat) ([]entities.Message, error) {
	publicChat := pc.db.Get("private chats")

	chats, ok := publicChat.(database.PrivateChatsData)
	if !ok {
		return nil, ErrDataError
	}

	return chats[chat], nil
}

func (pc *ChatsDB) GetUsersWithMessage(username string) ([]string, error) {
	publicChat := pc.db.Get("private chats")

	chats, ok := publicChat.(database.PrivateChatsData)
	if !ok {
		return nil, ErrDataError
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
