package private_message

import (
	"sync"

	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/internal/repository"
)

type PrivateChatsData map[entities.ChatMetadata][]entities.Message

type InMemoryDB interface {
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

func (pc *Repository) SendPrivateMessage(chat entities.ChatMetadata, msg entities.Message) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	privateChat := pc.db.Get(DBKey)

	chats, ok := privateChat.(PrivateChatsData)
	if !ok {
		return repository.ErrDataError
	}

	chats[chat] = append(chats[chat], msg)

	return nil
}

func (pc *Repository) GetPrivateChats() (map[entities.ChatMetadata][]entities.Message, error) {

	privateChat := pc.db.Get(DBKey)

	chats, ok := privateChat.(PrivateChatsData)
	if !ok {
		return nil, repository.ErrDataError
	}

	return chats, nil
}
