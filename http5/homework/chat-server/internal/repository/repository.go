package repository

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/database"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"
)

type InMemoryDB interface {
	Insert(key string, value interface{})
	Get(key string) interface{}
}

type Authorization interface {
	CreateUser(user entities.User) (string, int, error)
	GetUser(user entities.User) (entities.User, int, error)
}

type Chatting interface {
	SendPublicMessage(msg entities.Message) error
	SendPrivateMessage(chat database.Chat, msg entities.Message) error
	GetPublicMessages() ([]entities.Message, error)
	GetPrivateMessages(chat database.Chat) ([]entities.Message, error)
	GetUsersWithMessage(receiver string) ([]string, error)
}

type Repository struct {
	Auth Authorization
	Chat Chatting
}

func NewRepository(db InMemoryDB) *Repository {
	chatDB, ok := db.(*database.ChatDB)
	if !ok {
		return nil
	}

	return &Repository{
		Auth: NewAuthDB(chatDB),
		Chat: NewChattingDB(chatDB),
	}
}
