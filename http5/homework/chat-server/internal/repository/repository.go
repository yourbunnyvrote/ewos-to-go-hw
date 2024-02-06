package repository

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/chatutil"
)

type Authorization interface {
	CreateUser(user chatutil.User) (string, int, error)
	GetUser(user chatutil.User) (chatutil.User, int, error)
}

type Chatting interface {
	SendPublicMessage(msg chatutil.Message) error
	SendPrivateMessage(chat chatutil.Chat, msg chatutil.Message) error
	GetPublicMessages() ([]chatutil.Message, error)
	GetPrivateMessages(chat chatutil.Chat) ([]chatutil.Message, error)
	GetUsersWithMessage(receiver string) ([]string, error)
}

type Repository struct {
	Authorization
	Chatting
}

func NewRepository(db *chatutil.ChatDB) *Repository {
	return &Repository{
		Authorization: NewAuthDB(db),
		Chatting:      NewChatsDB(db),
	}
}
