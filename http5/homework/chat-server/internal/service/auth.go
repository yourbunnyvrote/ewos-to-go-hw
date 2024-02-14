package service

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"
)

type AuthService struct {
	repos repository.Authorization
}

func NewAuthService(repos repository.Authorization) *AuthService {
	return &AuthService{repos: repos}
}

func (as *AuthService) CreateUser(user entities.User) (string, error) {
	return as.repos.CreateUser(user)
}

func (as *AuthService) GetUser(username string) (entities.User, error) {
	return as.repos.GetUser(username)
}

func (as *AuthService) Identify(user entities.User) error {
	checkingUser, err := as.repos.GetUser(user.Username)
	if err != nil {
		return err
	}

	if user.Password != checkingUser.Password {
		return ErrIncorrectPassword
	}

	return nil
}
