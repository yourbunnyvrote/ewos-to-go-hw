package public_message

import "github.com/ew0s/ewos-to-go-hw/internal/domain/entities"

type PublicMessageRepo interface {
	SendPublicMessage(msg entities.Message) error
	GetPublicMessages() ([]entities.Message, error)
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

func (cs *Service) GetPublicMessages() ([]entities.Message, error) {
	return cs.repos.GetPublicMessages()
}

func (*Service) PaginateMessages(messages []entities.Message, params entities.PaginateParam) []entities.Message {
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
