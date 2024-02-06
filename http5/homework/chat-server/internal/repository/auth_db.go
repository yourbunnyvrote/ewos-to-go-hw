package repository

import (
	"errors"
	"net/http"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/chatutil"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type AuthDB struct {
	db *chatutil.ChatDB
}

func NewAuthDB(db *chatutil.ChatDB) *AuthDB {
	return &AuthDB{db: db}
}

func (a *AuthDB) CreateUser(user chatutil.User) (string, int, error) {
	if _, exists := a.db.Users[user.Username]; exists {
		return "", http.StatusConflict, ErrUserAlreadyExists
	}

	a.db.Users[user.Username] = user

	return user.Username, http.StatusCreated, nil
}

func (a *AuthDB) GetUser(user chatutil.User) (chatutil.User, int, error) {
	findUser, exist := a.db.Users[user.Username]
	if !exist {
		return chatutil.User{}, http.StatusUnprocessableEntity, ErrUserNotFound
	}

	return findUser, http.StatusOK, nil
}
