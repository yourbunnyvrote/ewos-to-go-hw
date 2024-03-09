package public_message

import (
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/pkg/pagination"
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

	pageMessages := pagination.Paginate(publicChat, params.Offset, params.Limit)

	return pageMessages, nil
}
