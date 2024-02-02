package repository

import (
	"errors"
	chatutil "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server"
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

func (a *AuthDB) GetUser(username, password string) (chatutil.User, error) {
	findUser, exist := a.db.Users[username]
	if !exist {
		return chatutil.User{}, ErrUserNotFound
	}

	if findUser.Password != password {
		return chatutil.User{}, ErrIncorrectPassword
	}

	return findUser, nil
}
