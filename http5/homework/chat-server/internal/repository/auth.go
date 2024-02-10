package repository

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/database"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"
)

type AuthDB struct {
	db *database.ChatDB
}

func NewAuthDB(db *database.ChatDB) *AuthDB {
	return &AuthDB{db: db}
}

func (a *AuthDB) CreateUser(user entities.User) (string, error) {
	err := a.db.AddUsersData(user)
	if err != nil {
		return "", err
	}

	return user.Username, nil
}

func (a *AuthDB) GetUser(username string) (entities.User, error) {
	findUser, err := a.db.GetUserData(username)
	if err != nil {
		return entities.User{}, err
	}

	return findUser, nil
}
