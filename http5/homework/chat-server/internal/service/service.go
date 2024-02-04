package service

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/chatutil"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"
)

type Authorization interface {
	CreateUser(user chatutil.User) (string, error)
	GetUser(user chatutil.User) (string, error)
}

type Chatting interface {
	SendPublicMessage(msg chatutil.Message) error
	SendPrivateMessage(chat chatutil.Chat, msg chatutil.Message) error
	GetPublicMessages() ([]chatutil.Message, error)
	GetPrivateMessages(chat chatutil.Chat) ([]chatutil.Message, error)
	GetUsersWithMessage(receiver string) ([]string, error)
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
