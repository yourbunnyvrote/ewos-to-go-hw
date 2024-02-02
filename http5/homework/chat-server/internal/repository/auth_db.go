package repository

import (
	"errors"
	chatutil "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type AuthDB struct {
	db *chatutil.ChatDB
}

func NewAuthDB(db *chatutil.ChatDB) *AuthDB {
	return &AuthDB{db: db}
}

func (as *AuthDB) CreateUser(user chatutil.User) (string, error) {
	if _, exists := as.db.Users[user.Username]; exists {
		return "", ErrUserAlreadyExists
	}

	as.db.Users[user.Username] = user

	return user.Username, nil
}
