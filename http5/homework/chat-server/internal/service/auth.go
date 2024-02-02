package service

import (
	chatutil "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"
)

type AuthService struct {
	repos repository.Authorization
}

func NewAuthService(repos repository.Authorization) *AuthService {
	return &AuthService{repos: repos}
}

func (as *AuthService) CreateUser(user chatutil.User) (string, error) {
	return as.repos.CreateUser(user)
}

func (as *AuthService) GetUser(username, password string) (string, error) {
	user, err := as.repos.GetUser(username, password)
	if err != nil {
		return "", err
	}

	return user.Username, nil
}
