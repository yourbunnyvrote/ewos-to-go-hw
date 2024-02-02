package service

import (
	chatutil "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"
)

type Authorization interface {
	CreateUser(user chatutil.User) (string, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{Authorization: NewAuthService(repos)}
}
