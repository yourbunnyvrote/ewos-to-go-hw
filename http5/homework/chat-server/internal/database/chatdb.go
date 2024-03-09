package database

import (
	"sync"
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
