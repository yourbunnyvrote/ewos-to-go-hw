package database

import (
	"sync"

	"github.com/ew0s/ewos-to-go-hw/internal/repository/inmemory/private_message"
	"github.com/ew0s/ewos-to-go-hw/internal/repository/inmemory/public_message"
	"github.com/ew0s/ewos-to-go-hw/internal/repository/inmemory/user"
)

type ChatDB struct {
	mu   *sync.RWMutex
	data map[string]interface{}
}

func NewChatDB() *ChatDB {
	db := ChatDB{
		mu:   &sync.RWMutex{},
		data: map[string]interface{}{},
	}

	db.data["users"] = user.UsersData{}
	db.data["public chats"] = public_message.PublicChatsData{}
	db.data["private chats"] = private_message.PrivateChatsData{}
	db.data["user list"] = private_message.PrivateMessageUserList{}

	return &db
}

func (db *ChatDB) Insert(key string, value interface{}) error {
	db.mu.Lock()
	db.data[key] = value
	db.mu.Unlock()

	return nil
}

func (db *ChatDB) Get(key string) (interface{}, error) {
	db.mu.RLock()
	data := db.data[key]
	db.mu.RUnlock()

	return data, nil
}
