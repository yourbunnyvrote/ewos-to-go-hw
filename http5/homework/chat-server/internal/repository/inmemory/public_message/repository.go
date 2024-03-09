package public_message

import (
	"sync"

	"github.com/ew0s/ewos-to-go-hw/internal/repository/inmemory"

	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
)

type PublicChatsData []entities.Message

type InMemoryDB interface {
	Insert(query string, _ interface{}) error
	Get(query string) (interface{}, error)
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

	publicChat, err := r.db.Get(DBKey)
	if err != nil {
		return err
	}

	messages, ok := publicChat.(PublicChatsData)
	if !ok {
		return inmemory.ErrDataError
	}

	messages = append(messages, msg)

	err = r.db.Insert(DBKey, messages)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetPublicChat() ([]entities.Message, error) {
	publicChat, err := r.db.Get(DBKey)
	if err != nil {
		return nil, err
	}

	messages, ok := publicChat.(PublicChatsData)
	if !ok {
		return nil, inmemory.ErrDataError
	}

	return messages, nil
}
