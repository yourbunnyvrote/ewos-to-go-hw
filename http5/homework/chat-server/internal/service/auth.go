package service

import (
	"errors"
	"net/http"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"
)

var (
	ErrEmptyCredentials  = errors.New("username or password is empty")
	ErrIncorrectPassword = errors.New("incorrect password")
)

type AuthService struct {
	repos repository.Authorization
}

func NewAuthService(repos repository.Authorization) *AuthService {
	return &AuthService{repos: repos}
}

func (as *AuthService) CreateUser(user entities.User) (string, int, error) {
	if user.Username == "" || user.Password == "" {
		return "", http.StatusBadRequest, ErrEmptyCredentials
	}

	return as.repos.CreateUser(user)
}

func (as *AuthService) GetUser(user entities.User) (string, int, error) {
	if user.Username == "" || user.Password == "" {
		return "", http.StatusBadRequest, ErrEmptyCredentials
	}

	gettingUser, statusCode, err := as.repos.GetUser(user)
	if err != nil {
		return "", statusCode, err
	}

	if gettingUser.Password != user.Password {
		return "", http.StatusUnauthorized, ErrIncorrectPassword
	}

	return gettingUser.Username, http.StatusOK, nil
}
