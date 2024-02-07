package repository

import (
	"errors"
	"net/http"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/database"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrDataError         = errors.New("incorrect data in memory db")
)

type AuthDB struct {
	db *database.ChatDB
}

func NewAuthDB(db *database.ChatDB) *AuthDB {
	return &AuthDB{db: db}
}

func (a *AuthDB) CreateUser(user entities.User) (string, int, error) {
	usersData := a.db.Get("users")

	users, ok := usersData.(database.UsersData)
	if !ok {
		return "", http.StatusInternalServerError, ErrDataError
	}

	if _, exists := users[user.Username]; exists {
		return "", http.StatusConflict, ErrUserAlreadyExists
	}

	users[user.Username] = user

	a.db.Insert("users", users)

	return user.Username, http.StatusCreated, nil
}

func (a *AuthDB) GetUser(user entities.User) (entities.User, int, error) {
	usersData := a.db.Get("users")

	users, ok := usersData.(database.UsersData)
	if !ok {
		return entities.User{}, http.StatusInternalServerError, ErrDataError
	}

	findUser, exist := users[user.Username]
	if !exist {
		return entities.User{}, http.StatusUnprocessableEntity, ErrUserNotFound
	}

	return findUser, http.StatusOK, nil
}
