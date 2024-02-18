package auth

import (
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/internal/service"
)

type Authorization interface {
	CreateUser(user entities.AuthCredentials) error
	GetUser(username string) (entities.AuthCredentials, error)
}

type Service struct {
	repos Authorization
}

func NewService(repos Authorization) *Service {
	return &Service{
		repos: repos,
	}
}

func (as *Service) CreateUser(credentials entities.AuthCredentials) error {
	return as.repos.CreateUser(credentials)
}

func (as *Service) GetUser(username string) (entities.AuthCredentials, error) {
	return as.repos.GetUser(username)
}

func (as *Service) Identify(user entities.AuthCredentials) error {
	checkingUser, err := as.repos.GetUser(user.Login)
	if err != nil {
		return err
	}

	if user.Password != checkingUser.Password {
		return service.ErrIncorrectPassword
	}

	return nil
}
