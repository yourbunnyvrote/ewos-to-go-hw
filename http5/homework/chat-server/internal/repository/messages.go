package repository

import chatutil "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server"

type PublicChatDB struct {
	db *chatutil.ChatDB
}

func NewPublicChatDB(db *chatutil.ChatDB) *PublicChatDB {
	return &PublicChatDB{db: db}
}

func (pc *PublicChatDB) SendMessage(msg chatutil.Message) error {
	pc.db.PublicChat = append(pc.db.PublicChat, msg)
	return nil
}

func (pc *PublicChatDB) GetMessage() ([]chatutil.Message, error) {
	return pc.db.PublicChat, nil
}
