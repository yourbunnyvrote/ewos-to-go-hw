package repository

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"
)

type InMemoryDBChat interface {
	AddPublicMessage(msg entities.Message) error
	GetPublicMessages() ([]entities.Message, error)
	AddPrivateMessage(chat entities.Chat, msg entities.Message) error
	GetPrivateMessages(chat entities.Chat) ([]entities.Message, error)
	GetUsersPrivateMessages(username string) ([]string, error)
}

type ChattingDB struct {
	db InMemoryDBChat
}

func NewChattingDB(db InMemoryDBChat) *ChattingDB {
	return &ChattingDB{db: db}
}

func (pc *ChattingDB) SendPublicMessage(msg entities.Message) error {
	return pc.db.AddPublicMessage(msg)
}

func (pc *ChattingDB) SendPrivateMessage(chat entities.Chat, msg entities.Message) error {
	return pc.db.AddPrivateMessage(chat, msg)
}

func (pc *ChattingDB) GetPublicMessages() ([]entities.Message, error) {
	return pc.db.GetPublicMessages()
}

func (pc *ChattingDB) GetPrivateMessages(chat entities.Chat) ([]entities.Message, error) {
	return pc.db.GetPrivateMessages(chat)
}

func (pc *ChattingDB) GetUsersWithMessage(username string) ([]string, error) {
	return pc.db.GetUsersPrivateMessages(username)
}
