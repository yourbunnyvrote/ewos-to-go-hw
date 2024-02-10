package database

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/constants"
)

func (db *ChatDB) AddUsersData(user entities.User) error {
	if user.Username == "" || user.Password == "" {
		return constants.ErrEmptyCredentials
	}

	usersData := db.Get("users")

	users, ok := usersData.(UsersData)
	if !ok {
		return constants.ErrDataError
	}

	if _, exists := users[user.Username]; exists {
		return constants.ErrUserAlreadyExists
	}

	users[user.Username] = user

	return nil
}

func (db *ChatDB) GetUserData(username string) (entities.User, error) {
	if username == "" {
		return entities.User{}, constants.ErrBadRequest
	}

	usersData := db.Get("users")

	users, ok := usersData.(UsersData)
	if !ok {
		return entities.User{}, constants.ErrDataError
	}

	findUser, exist := users[username]
	if !exist {
		return entities.User{}, constants.ErrUserNotFound
	}

	return findUser, nil
}
