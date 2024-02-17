package repository

import (
	"sync"

	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
)

type ChattingDB struct {
	mu *sync.RWMutex
	db InMemoryDB
}

func NewChattingDB(db InMemoryDB) *ChattingDB {
	return &ChattingDB{
		mu: &sync.RWMutex{},
		db: db,
	}
}

func (pc *ChattingDB) SendPublicMessage(msg entities.Message) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	publicChat := pc.db.Get("public chats")

	messages, ok := publicChat.(PublicChatsData)
	if !ok {
		return ErrorDataError
	}

	messages = append(messages, msg)

	pc.db.Insert("public chats", messages)

	return nil
}

func (pc *ChattingDB) SendPrivateMessage(chat entities.UsersPair, msg entities.Message) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	publicChat := pc.db.Get("private chats")

	chats, ok := publicChat.(PrivateChatsData)
	if !ok {
		return ErrorDataError
	}

	if chat.User1 < chat.User2 {
		chat.User1, chat.User2 = chat.User2, chat.User1
	}

	chats[chat] = append(chats[chat], msg)

	return nil
}

func (pc *ChattingDB) GetPublicMessages() ([]entities.Message, error) {
	publicChat := pc.db.Get("public chats")

	messages, ok := publicChat.(PublicChatsData)
	if !ok {
		return nil, ErrorDataError
	}

	return messages, nil
}

func (pc *ChattingDB) GetPrivateMessages(chat entities.UsersPair) ([]entities.Message, error) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	privateChat := pc.db.Get("private chats")

	chats, ok := privateChat.(PrivateChatsData)
	if !ok {
		return nil, ErrorDataError
	}

	if chat.User1 < chat.User2 {
		chat.User1, chat.User2 = chat.User2, chat.User1
	}

	messages := chats[chat]

	return messages, nil
}

func (pc *ChattingDB) GetUsersWithMessage(username string) ([]string, error) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	publicChat := pc.db.Get("private chats")

	chats, ok := publicChat.(PrivateChatsData)
	if !ok {
		return nil, ErrorDataError
	}

	listUsers := make([]string, 0)

	for key := range chats {
		if key.User1 == username || key.User2 == username {
			if key.User1 == username {
				listUsers = append(listUsers, key.User2)
			} else {
				listUsers = append(listUsers, key.User1)
			}
		}
	}

	return listUsers, nil
}
