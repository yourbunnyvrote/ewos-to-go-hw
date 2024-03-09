package private_message

import (
	"sync"

	"github.com/ew0s/ewos-to-go-hw/internal/api/mapper"
	"github.com/ew0s/ewos-to-go-hw/internal/repository/inmemory"

	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
)

type (
	PrivateChatsData       map[entities.ChatMetadata][]entities.Message
	PrivateMessageUserList map[string]map[string]struct{}
)

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

func (pc *Repository) SendPrivateMessage(receiver string, msg entities.Message) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	privateChat, err := pc.db.Get(DBKeyChats)
	if err != nil {
		return err
	}

	chats, ok := privateChat.(PrivateChatsData)
	if !ok {
		return inmemory.ErrDataError
	}

	chat := mapper.MakeChatMetadata(receiver, msg.Username)

	if _, exist := chats[chat]; !exist {
		userList, err := pc.db.Get(DBKeyUserList)
		if err != nil {
			return err
		}

		users, ok := userList.(PrivateMessageUserList)
		if !ok {
			return inmemory.ErrDataError
		}

		users[chat.Username1][chat.Username2] = struct{}{}
		users[chat.Username2][chat.Username1] = struct{}{}
	}

	chats[chat] = append(chats[chat], msg)

	return nil
}

func (pc *Repository) GetPrivateChat(chat entities.ChatMetadata) ([]entities.Message, error) {
	privateChat, err := pc.db.Get(DBKeyChats)
	if err != nil {
		return nil, err
	}

	chats, ok := privateChat.(PrivateChatsData)
	if !ok {
		return nil, inmemory.ErrDataError
	}

	return chats[chat], nil
}

func (pc *Repository) GetUserList(username string) ([]string, error) {
	userList, err := pc.db.Get(DBKeyUserList)
	if err != nil {
		return nil, err
	}

	users, ok := userList.(PrivateMessageUserList)
	if !ok {
		return nil, inmemory.ErrDataError
	}

	result := make([]string, 0)
	for key := range users[username] {
		result = append(result, key)
	}

	return result, nil
}
