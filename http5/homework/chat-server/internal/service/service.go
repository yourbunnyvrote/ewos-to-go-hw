package service

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"
)

type Authorization interface {
	CreateUser(user entities.User) (string, error)
	GetUser(username string) (entities.User, error)
}

type Chatting interface {
	SendPublicMessage(msg entities.Message) error
	SendPrivateMessage(chat entities.Chat, msg entities.Message) error
	GetPublicMessages() ([]entities.Message, error)
	GetPrivateMessages(chat entities.Chat) ([]entities.Message, error)
	GetUsersWithMessage(receiver string) ([]string, error)
}

type Service struct {
	Auth Authorization
	Chat Chatting
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Auth: NewAuthService(repos.Auth),
		Chat: NewChatService(repos.Chat),
	}
}
