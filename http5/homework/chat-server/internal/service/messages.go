package service

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/chatutil"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"
)

type ChatService struct {
	repos repository.Chatting
}

func NewChatService(repos repository.Chatting) *ChatService {
	return &ChatService{repos: repos}
}

func (cs *ChatService) SendPublicMessage(msg chatutil.Message) error {
	return cs.repos.SendPublicMessage(msg)
}

func (cs *ChatService) SendPrivateMessage(chat chatutil.Chat, msg chatutil.Message) error {
	if chat.User1 < chat.User2 {
		chat.User1, chat.User2 = chat.User2, chat.User1
	}

	return cs.repos.SendPrivateMessage(chat, msg)
}

func (cs *ChatService) GetPublicMessages() ([]chatutil.Message, error) {
	return cs.repos.GetPublicMessages()
}

func (cs *ChatService) GetPrivateMessages(chat chatutil.Chat) ([]chatutil.Message, error) {
	if chat.User1 < chat.User2 {
		chat.User1, chat.User2 = chat.User2, chat.User1
	}

	return cs.repos.GetPrivateMessages(chat)
}

func (cs *ChatService) GetUsersWithMessage(receiver string) ([]string, error) {
	return cs.repos.GetUsersWithMessage(receiver)
}
