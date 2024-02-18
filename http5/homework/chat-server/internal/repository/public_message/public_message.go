package public_message

import (
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/internal/repository"
	"sync"
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

func (r *Repository) SendPublicMessage(msg entities.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	publicChat := r.db.Get("public chats")

	messages, ok := publicChat.(repository.PublicChatsData)
	if !ok {
		return repository.ErrDataError
	}

	messages = append(messages, msg)

	r.db.Insert("public chats", messages)

	return nil
}

func (r *Repository) GetPublicMessages() ([]entities.Message, error) {
	publicChat := r.db.Get("public chats")

	messages, ok := publicChat.(repository.PublicChatsData)
	if !ok {
		return nil, repository.ErrDataError
	}

	return messages, nil
}
