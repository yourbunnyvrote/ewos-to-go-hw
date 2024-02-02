package service

import (
	chatutil "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"
)

type Authorization interface {
	CreateUser(user chatutil.User) (string, error)
	GetUser(username, password string) (string, error)
}

type Chatting interface {
	SendMessage(msg chatutil.Message) error
	GetMessage() ([]chatutil.Message, error)
}

type Service struct {
	Authorization
	Chatting
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Chatting:      NewChatService(repos.Chatting),
	}
}
