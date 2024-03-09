package auth

import (
	"errors"

	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/auth/request"
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/auth/response"
	"github.com/ew0s/ewos-to-go-hw/internal/api/mapper"
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/internal/service"
)

var ErrIncorrectTypeUser = errors.New("user are not of type entities.AuthCredentials")

type BasicAuthRepo interface {
	CreateUser(user entities.AuthCredentials) error
	GetUser(username string) (entities.AuthCredentials, error)
}

type BasicAuthService struct {
	repos BasicAuthRepo
}

func NewBasicAuthService(repos BasicAuthRepo) *BasicAuthService {
	return &BasicAuthService{
		repos: repos,
	}
}

func (s *BasicAuthService) CreateUser(req request.RegistrationRequest) (interface{}, error) {
	creds := mapper.MakeEntityAuthCredentials(req.Username, req.Password)

	err := s.repos.CreateUser(creds)
	if err != nil {
		return response.RegistrationResponse{}, err
	}

	return mapper.MakeRegistrationResponse(creds), nil
}

func (s *BasicAuthService) GetUser(username string) (entities.AuthCredentials, error) {
	return s.repos.GetUser(username)
}

func (s *BasicAuthService) Identify(creds interface{}) (string, error) {
	user, ok := creds.(entities.AuthCredentials)
	if !ok {
		return "", ErrIncorrectTypeUser
	}

	checkingUser, err := s.repos.GetUser(user.Login)
	if err != nil {
		return "", err
	}

	if user.Password != checkingUser.Password {
		return "", service.ErrIncorrectPassword
	}

	return user.Login, nil
}
