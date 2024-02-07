package database

import "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"

type ChatDB struct {
	data map[string]interface{}
}

func NewChatDB() *ChatDB {
	return &ChatDB{
		data: map[string]interface{}{
			"users":         UsersData{},
			"public chats":  PublicChatsData{},
			"private chats": PrivateChatsData{},
		},
	}
}

func (c *ChatDB) Insert(key string, value interface{}) {
	c.data[key] = value
}

func (c *ChatDB) Get(key string) interface{} {
	return c.data[key]
}

type (
	UsersData        map[string]entities.User
	PublicChatsData  []entities.Message
	PrivateChatsData map[Chat][]entities.Message
)

type Chat struct {
	User1 string
	User2 string
}
