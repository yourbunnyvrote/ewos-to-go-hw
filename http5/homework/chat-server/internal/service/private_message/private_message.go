package private_message

import "github.com/ew0s/ewos-to-go-hw/internal/domain/entities"

type PrivateMessageRepo interface {
	SendPrivateMessage(chat entities.ChatMetadata, msg entities.Message) error
	GetPrivateMessages(chat entities.ChatMetadata) ([]entities.Message, error)
	GetUsersWithMessage(receiver string) ([]string, error)
}

type Service struct {
	repos PrivateMessageRepo
}

func NewService(repos PrivateMessageRepo) *Service {
	return &Service{repos: repos}
}

func (cs *Service) SendPrivateMessage(chat entities.ChatMetadata, msg entities.Message) error {
	return cs.repos.SendPrivateMessage(chat, msg)
}

func (cs *Service) GetPrivateMessages(chat entities.ChatMetadata) ([]entities.Message, error) {
	return cs.repos.GetPrivateMessages(chat)
}

func (cs *Service) GetUsersWithMessage(receiver string) ([]string, error) {
	return cs.repos.GetUsersWithMessage(receiver)
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
