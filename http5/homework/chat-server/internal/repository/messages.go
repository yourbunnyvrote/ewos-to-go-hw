package repository

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/chatutil"
)

type ChatsDB struct {
	db *chatutil.ChatDB
}

func NewChatsDB(db *chatutil.ChatDB) *ChatsDB {
	return &ChatsDB{db: db}
}

func (pc *ChatsDB) SendPublicMessage(msg chatutil.Message) error {
	pc.db.PublicChat = append(pc.db.PublicChat, msg)
	return nil
}

func (pc *ChatsDB) SendPrivateMessage(chat chatutil.Chat, msg chatutil.Message) error {
	pc.db.PrivateChats[chat] = append(pc.db.PrivateChats[chat], msg)
	return nil
}

func (pc *ChatsDB) GetPublicMessages() ([]chatutil.Message, error) {
	return pc.db.PublicChat, nil
}

func (pc *ChatsDB) GetPrivateMessages(chat chatutil.Chat) ([]chatutil.Message, error) {
	return pc.db.PrivateChats[chat], nil
}

func (pc *ChatsDB) GetUsersWithMessage(username string) ([]string, error) {
	result := make([]string, 0)

	for key := range pc.db.PrivateChats {
		if key.User1 == username || key.User2 == username {
			if key.User1 == username {
				result = append(result, key.User2)
			} else {
				result = append(result, key.User1)
			}
		}
	}

	return result, nil
}
