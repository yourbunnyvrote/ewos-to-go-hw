package service

import (
	chatutil "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"
)

type ChatService struct {
	repos repository.Chatting
}

func NewChatService(repos repository.Chatting) *ChatService {
	return &ChatService{repos: repos}
}

func (cs *ChatService) SendMessage(msg chatutil.Message) error {
	return cs.repos.SendMessage(msg)
}

func (cs *ChatService) GetMessage() ([]chatutil.Message, error) {
	return cs.repos.GetMessage()
}
