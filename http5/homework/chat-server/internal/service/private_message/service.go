package private_message

import (
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/pkg/pagination"
)

type PrivateMessageRepo interface {
	SendPrivateMessage(receiver string, msg entities.Message) error
	GetPrivateChat(chat entities.ChatMetadata) ([]entities.Message, error)
	GetUserList(username string) ([]string, error)
}

type Service struct {
	repos PrivateMessageRepo
}

func NewService(repos PrivateMessageRepo) *Service {
	return &Service{repos: repos}
}

func (cs *Service) SendPrivateMessage(receiver string, msg entities.Message) error {
	return cs.repos.SendPrivateMessage(receiver, msg)
}

func (cs *Service) GetPrivateMessages(chat entities.ChatMetadata, params entities.PaginateParam) ([]entities.Message, error) {
	messages, err := cs.repos.GetPrivateChat(chat)
	if err != nil {
		return nil, err
	}

	pageMessages := pagination.Paginate(messages, params.Offset, params.Limit)

	return pageMessages, nil
}

func (cs *Service) GetUsersWithMessage(receiver string) ([]string, error) {
	return cs.repos.GetUserList(receiver)
}
