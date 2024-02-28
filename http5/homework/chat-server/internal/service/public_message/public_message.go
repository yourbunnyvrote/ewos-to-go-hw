package public_message

import (
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/internal/service/message"
)

type PublicMessageRepo interface {
	SendPublicMessage(msg entities.Message) error
	GetPublicChat() ([]entities.Message, error)
}

type Service struct {
	repos PublicMessageRepo
}

func NewService(repos PublicMessageRepo) *Service {
	return &Service{repos: repos}
}

func (cs *Service) SendPublicMessage(msg entities.Message) error {
	return cs.repos.SendPublicMessage(msg)
}

func (cs *Service) GetPublicMessages(params entities.PaginateParam) ([]entities.Message, error) {
	publicChat, err := cs.repos.GetPublicChat()
	if err != nil {
		return nil, err
	}

	pageMessages := message.PaginateMessages(publicChat, params)

	return pageMessages, nil
}
