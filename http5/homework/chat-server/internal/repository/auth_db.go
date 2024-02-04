package repository

import (
	"errors"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/chatutil"
)

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrIncorrectPassword = errors.New("incorrect password")

type AuthDB struct {
	db *chatutil.ChatDB
}

func NewAuthDB(db *chatutil.ChatDB) *AuthDB {
	return &AuthDB{db: db}
}

func (a *AuthDB) CreateUser(user chatutil.User) (string, error) {
	if _, exists := a.db.Users[user.Username]; exists {
		return "", ErrUserAlreadyExists
	}

	a.db.Users[user.Username] = user

	return user.Username, nil
}

func (a *AuthDB) GetUser(user chatutil.User) (chatutil.User, error) {
	findUser, exist := a.db.Users[user.Username]
	if !exist {
		return chatutil.User{}, ErrUserNotFound
	}

	if findUser.Password != user.Password {
		return chatutil.User{}, ErrIncorrectPassword
	}

	return findUser, nil
}
