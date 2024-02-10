package database

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"
	"sync"
)

type ChatDB struct {
	mu   *sync.RWMutex
	data map[string]interface{}
}

type (
	UsersData        map[string]entities.User
	PublicChatsData  []entities.Message
	PrivateChatsData map[entities.Chat][]entities.Message
)

func NewChatDB() *ChatDB {
	db := ChatDB{
		mu:   &sync.RWMutex{},
		data: map[string]interface{}{},
	}

	db.Insert("users", UsersData{})
	db.Insert("public chats", PublicChatsData{})
	db.Insert("private chats", PrivateChatsData{})

	return &db
}

func (db *ChatDB) Insert(key string, value interface{}) {
	db.mu.Lock()
	db.data[key] = value
	db.mu.Unlock()
}

func (db *ChatDB) Get(key string) interface{} {
	db.mu.RLock()
	data := db.data[key]
	db.mu.RUnlock()
	return data
}
