package user

import (
	"sync"

	"github.com/ew0s/ewos-to-go-hw/internal/repository/inmemory"

	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
)

type UsersData map[string]entities.AuthCredentials

type InMemoryDB interface {
	Get(query string) (data interface{}, _ error)
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

	usersData, err := a.db.Get(DBKey)
	if err != nil {
		return err
	}

	users, ok := usersData.(UsersData)
	if !ok {
		return inmemory.ErrDataError
	}

	if _, exists := users[credentials.Login]; exists {
		return ErrUserAlreadyExists
	}

	users[credentials.Login] = credentials

	return nil
}

func (a *Repository) GetUser(username string) (entities.AuthCredentials, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	usersData, err := a.db.Get(DBKey)
	if err != nil {
		return entities.AuthCredentials{}, err
	}

	users, ok := usersData.(UsersData)
	if !ok {
		return entities.AuthCredentials{}, inmemory.ErrDataError
	}

	findUser, exist := users[username]

	if !exist {
		return entities.AuthCredentials{}, ErrUserNotFound
	}

	return findUser, nil
}
