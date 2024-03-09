package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/auth/request"
	"github.com/ew0s/ewos-to-go-hw/internal/api/mapper"
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
)

type JWTAuthRepo interface {
	CreateUser(user entities.AuthCredentials) error
	GetUser(username string) (entities.AuthCredentials, error)
}

const (
	signingKey          = "-aenL9S^hXG18FXt4KdMvh8ozbZNkl"
	tokenAvailableHours = 12
)

type tokenClaims struct {
	jwt.StandardClaims
	Login string
}

type JWTService struct {
	repos JWTAuthRepo
}

func NewJWTService(repos JWTAuthRepo) *JWTService {
	return &JWTService{
		repos: repos,
	}
}

func (s *JWTService) CreateUser(req request.RegistrationRequest) (interface{}, error) {
	creds := mapper.MakeEntityAuthCredentials(req.Username, req.Password)

	err := s.repos.CreateUser(creds)
	if err != nil {
		return struct{}{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenAvailableHours * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		creds.Login,
	})

	tokenStr, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return struct{}{}, err
	}

	resp := mapper.MakeJWTResponse(tokenStr)

	return resp, nil
}

func (s *JWTService) GetUser(username string) (entities.AuthCredentials, error) {
	return s.repos.GetUser(username)
}

func (s *JWTService) Identify(accessToken interface{}) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken.(string), &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", ErrIncorrectTypeClaims
	}

	user, err := s.repos.GetUser(claims.Login)
	if err != nil {
		return "", err
	}

	return user.Login, nil
}
