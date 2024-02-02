package repository

import chatutil "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server"

type Authorization interface {
	CreateUser(user chatutil.User) (string, error)
	GetUser(username, password string) (chatutil.User, error)
}

type Chatting interface {
	SendMessage(msg chatutil.Message) error
	GetMessage() ([]chatutil.Message, error)
}

type Repository struct {
	Authorization
	Chatting
}

func NewRepository(db *chatutil.ChatDB) *Repository {
	return &Repository{
		Authorization: NewAuthDB(db),
		Chatting:      NewPublicChatDB(db),
	}
}
