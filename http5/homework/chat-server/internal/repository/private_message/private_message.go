package private_message

import (
	"sync"

	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/internal/repository"
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

func (pc *Repository) SendPrivateMessage(chat entities.ChatMetadata, msg entities.Message) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	publicChat := pc.db.Get("private chats")

	chats, ok := publicChat.(repository.PrivateChatsData)
	if !ok {
		return repository.ErrDataError
	}

	if chat.Username1 < chat.Username2 {
		chat.Username1, chat.Username2 = chat.Username2, chat.Username1
	}

	chats[chat] = append(chats[chat], msg)

	return nil
}

func (pc *Repository) GetPrivateMessages(chat entities.ChatMetadata) ([]entities.Message, error) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	privateChat := pc.db.Get("private chats")

	chats, ok := privateChat.(repository.PrivateChatsData)
	if !ok {
		return nil, repository.ErrDataError
	}

	if chat.Username1 < chat.Username2 {
		chat.Username1, chat.Username2 = chat.Username2, chat.Username1
	}

	messages := chats[chat]

	return messages, nil
}

func (pc *Repository) GetUsersWithMessage(username string) ([]string, error) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	publicChat := pc.db.Get("private chats")

	chats, ok := publicChat.(repository.PrivateChatsData)
	if !ok {
		return nil, repository.ErrDataError
	}

	listUsers := make([]string, 0)

	for key := range chats {
		if key.Username1 == username || key.Username2 == username {
			if key.Username1 == username {
				listUsers = append(listUsers, key.Username2)
			} else {
				listUsers = append(listUsers, key.Username1)
			}
		}
	}

	return listUsers, nil
}
