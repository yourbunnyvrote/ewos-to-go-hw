package service

import (
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/internal/repository"
)

type ChatService struct {
	repos repository.Chatting
}

func NewChatService(repos repository.Chatting) *ChatService {
	return &ChatService{repos: repos}
}

func (cs *ChatService) SendPublicMessage(msg entities.Message) error {
	return cs.repos.SendPublicMessage(msg)
}

func (cs *ChatService) SendPrivateMessage(chat entities.UsersPair, msg entities.Message) error {
	return cs.repos.SendPrivateMessage(chat, msg)
}

func (cs *ChatService) GetPublicMessages() ([]entities.Message, error) {
	return cs.repos.GetPublicMessages()
}

func (cs *ChatService) GetPrivateMessages(chat entities.UsersPair) ([]entities.Message, error) {
	return cs.repos.GetPrivateMessages(chat)
}

func (cs *ChatService) GetUsersWithMessage(receiver string) ([]string, error) {
	return cs.repos.GetUsersWithMessage(receiver)
}

func (*ChatService) PaginateMessages(messages []entities.Message, params entities.PaginateParam) []entities.Message {
	startIndex := params.Offset
	endIndex := startIndex + params.Limit

	if endIndex > len(messages) {
		endIndex = len(messages)
	}

	if startIndex > len(messages) {
		startIndex = len(messages)
	}

	return messages[startIndex:endIndex]
}
