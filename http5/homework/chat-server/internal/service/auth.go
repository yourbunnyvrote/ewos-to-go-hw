package service

import (
	"errors"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/chatutil"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"
)

var ErrEmptyCredentials = errors.New("username or password is empty")

type AuthService struct {
	repos repository.Authorization
}

func NewAuthService(repos repository.Authorization) *AuthService {
	return &AuthService{repos: repos}
}

func (as *AuthService) CreateUser(user chatutil.User) (string, error) {
	if user.Username == "" || user.Password == "" {
		return "", ErrEmptyCredentials
	}

	return as.repos.CreateUser(user)
}

func (as *AuthService) GetUser(user chatutil.User) (string, error) {
	if user.Username == "" || user.Password == "" {
		return "", ErrEmptyCredentials
	}

	user, err := as.repos.GetUser(user)
	if err != nil {
		return "", err
	}

	return user.Username, nil
}
