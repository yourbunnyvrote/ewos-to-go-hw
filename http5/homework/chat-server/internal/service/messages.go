package service

import (
	"errors"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/database"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"
)

var ErrMsgEmpty = errors.New("message is empty")

type ChatService struct {
	repos repository.Chatting
}

func NewChatService(repos repository.Chatting) *ChatService {
	return &ChatService{repos: repos}
}

func (cs *ChatService) SendPublicMessage(msg entities.Message) error {
	if msg.Content == "" {
		return ErrMsgEmpty
	}

	return cs.repos.SendPublicMessage(msg)
}

func (cs *ChatService) SendPrivateMessage(chat database.Chat, msg entities.Message) error {
	if msg.Content == "" {
		return ErrMsgEmpty
	}

	if chat.User1 < chat.User2 {
		chat.User1, chat.User2 = chat.User2, chat.User1
	}

	return cs.repos.SendPrivateMessage(chat, msg)
}

func (cs *ChatService) GetPublicMessages() ([]entities.Message, error) {
	return cs.repos.GetPublicMessages()
}

func (cs *ChatService) GetPrivateMessages(chat database.Chat) ([]entities.Message, error) {
	if chat.User1 < chat.User2 {
		chat.User1, chat.User2 = chat.User2, chat.User1
	}

	return cs.repos.GetPrivateMessages(chat)
}

func (cs *ChatService) GetUsersWithMessage(receiver string) ([]string, error) {
	return cs.repos.GetUsersWithMessage(receiver)
}
