package private_message

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

func (pc *Repository) SendPrivateMessage(chat entities.UsersPair, msg entities.Message) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	publicChat := pc.db.Get("private chats")

	chats, ok := publicChat.(repository.PrivateChatsData)
	if !ok {
		return repository.ErrorDataError
	}

	if chat.User1 < chat.User2 {
		chat.User1, chat.User2 = chat.User2, chat.User1
	}

	chats[chat] = append(chats[chat], msg)

	return nil
}

func (pc *Repository) GetPrivateMessages(chat entities.UsersPair) ([]entities.Message, error) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	privateChat := pc.db.Get("private chats")

	chats, ok := privateChat.(repository.PrivateChatsData)
	if !ok {
		return nil, repository.ErrorDataError
	}

	if chat.User1 < chat.User2 {
		chat.User1, chat.User2 = chat.User2, chat.User1
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
		return nil, repository.ErrorDataError
	}

	listUsers := make([]string, 0)

	for key := range chats {
		if key.User1 == username || key.User2 == username {
			if key.User1 == username {
				listUsers = append(listUsers, key.User2)
			} else {
				listUsers = append(listUsers, key.User1)
			}
		}
	}

	return listUsers, nil
}
