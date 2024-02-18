package auth

import (
	"github.com/ew0s/ewos-to-go-hw/internal/repository"
	"sync"

	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
)

type InMemoryDB interface {
	Insert(key string, value interface{})
	Get(key string) interface{}
}

type Repository struct {
	mu *sync.RWMutex
	db InMemoryDB
}

func NewRepository(db InMemoryDB) *Repository {
	return &Repository{
		mu: &sync.RWMutex{},
		db: db,
	}
}

func (a *Repository) CreateUser(user entities.User) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	usersData := a.db.Get("users")

	users, ok := usersData.(repository.UsersData)
	if !ok {
		return "", repository.ErrorDataError
	}

	_, exists := users[user.Username]

	if exists {
		return "", repository.ErrorUserAlreadyExists
	}

	users[user.Username] = user

	return user.Username, nil
}

func (a *Repository) GetUser(username string) (entities.User, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	usersData := a.db.Get("users")

	users, ok := usersData.(repository.UsersData)
	if !ok {
		return entities.User{}, repository.ErrorDataError
	}

	findUser, exist := users[username]

	if !exist {
		return entities.User{}, repository.ErrorUserNotFound
	}

	return findUser, nil
}
