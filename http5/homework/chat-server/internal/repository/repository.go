package repository

import (
	"github.com/ew0s/ewos-to-go-hw/internal/database"
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
)

type InMemoryDB interface {
	Insert(key string, value interface{})
	Get(key string) interface{}
}

type Authorization interface {
	CreateUser(user entities.User) (string, error)
	GetUser(username string) (entities.User, error)
}

type Chatting interface {
	SendPublicMessage(msg entities.Message) error
	SendPrivateMessage(chat entities.UsersPair, msg entities.Message) error
	GetPublicMessages() ([]entities.Message, error)
	GetPrivateMessages(chat entities.UsersPair) ([]entities.Message, error)
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

	chatDB.Insert("users", UsersData{})
	chatDB.Insert("public chats", PublicChatsData{})
	chatDB.Insert("private chats", PrivateChatsData{})

	return &Repository{
		Auth: NewAuthDB(chatDB),
		Chat: NewChattingDB(chatDB),
	}
}
