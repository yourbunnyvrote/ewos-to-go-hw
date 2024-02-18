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

func (a *Repository) CreateUser(credentials entities.AuthCredentials) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	usersData := a.db.Get("users")

	users, ok := usersData.(repository.UsersData)
	if !ok {
		return repository.ErrDataError
	}

	if _, exists := users[credentials.Login]; exists {
		return repository.ErrUserAlreadyExists
	}

	users[credentials.Login] = credentials

	return nil
}

func (a *Repository) GetUser(username string) (entities.AuthCredentials, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	usersData := a.db.Get("users")

	users, ok := usersData.(repository.UsersData)
	if !ok {
		return entities.AuthCredentials{}, repository.ErrDataError
	}

	findUser, exist := users[username]

	if !exist {
		return entities.AuthCredentials{}, repository.ErrUserNotFound
	}

	return findUser, nil
}
