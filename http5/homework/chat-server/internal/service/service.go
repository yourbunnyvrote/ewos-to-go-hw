package service

import (
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/internal/repository"
)

type Authorization interface {
	CreateUser(user entities.User) (string, error)
	GetUser(username string) (entities.User, error)
	Identify(user entities.AuthCredentials) error
}

type Chatting interface {
	SendPublicMessage(msg entities.Message) error
	SendPrivateMessage(chat entities.UsersPair, msg entities.Message) error
	GetPublicMessages() ([]entities.Message, error)
	GetPrivateMessages(chat entities.UsersPair) ([]entities.Message, error)
	GetUsersWithMessage(receiver string) ([]string, error)
	PaginateMessages(messages []entities.Message, params entities.PaginateParam) []entities.Message
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
