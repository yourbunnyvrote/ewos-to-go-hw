package database

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"
)

func (db *ChatDB) AddUsersData(user entities.User) error {
	if user.Username == "" {
		return ErrorUsernameEmpty
	} else if user.Password == "" {
		return ErrorPasswordEmpty
	}

	usersData := db.Get("users")

	users, ok := usersData.(UsersData)
	if !ok {
		return ErrorDataError
	}

	db.mu.RLock()
	_, exists := users[user.Username]
	db.mu.RUnlock()

	if exists {
		return ErrorUserAlreadyExists
	}

	db.mu.Lock()
	users[user.Username] = user
	db.mu.RUnlock()

	return nil
}

func (db *ChatDB) GetUserData(username string) (entities.User, error) {
	if username == "" {
		return entities.User{}, ErrorUsernameEmpty
	}

	usersData := db.Get("users")

	users, ok := usersData.(UsersData)
	if !ok {
		return entities.User{}, ErrorDataError
	}

	db.mu.RLock()
	findUser, exist := users[username]
	db.mu.RUnlock()

	if !exist {
		return entities.User{}, ErrorUserNotFound
	}

	return findUser, nil
}
