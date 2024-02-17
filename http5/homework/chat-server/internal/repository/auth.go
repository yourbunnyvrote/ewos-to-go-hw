package repository

import (
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
	"sync"
)

type AuthDB struct {
	mu *sync.RWMutex
	db InMemoryDB
}

func NewAuthDB(db InMemoryDB) *AuthDB {
	return &AuthDB{
		mu: &sync.RWMutex{},
		db: db,
	}
}

func (a *AuthDB) CreateUser(user entities.User) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if user.Username == "" {
		return "", ErrorUsernameEmpty
	} else if user.Password == "" {
		return "", ErrorPasswordEmpty
	}

	usersData := a.db.Get("users")

	users, ok := usersData.(UsersData)
	if !ok {
		return "", ErrorDataError
	}

	_, exists := users[user.Username]

	if exists {
		return "", ErrorUserAlreadyExists
	}

	users[user.Username] = user

	return user.Username, nil
}

func (a *AuthDB) GetUser(username string) (entities.User, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if username == "" {
		return entities.User{}, ErrorUsernameEmpty
	}

	usersData := a.db.Get("users")

	users, ok := usersData.(UsersData)
	if !ok {
		return entities.User{}, ErrorDataError
	}

	findUser, exist := users[username]

	if !exist {
		return entities.User{}, ErrorUserNotFound
	}

	return findUser, nil
}
