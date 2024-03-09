package private_message

import (
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/pkg/pagination"
)

type PrivateMessagePostgresRepo interface {
	SendPrivateMessage(receiver string, msg entities.Message) error
	GetPrivateChat(chat entities.ChatMetadata) ([]entities.Message, error)
	GetUserList(username string) ([]string, error)
}

type PostgresService struct {
	repos PrivateMessagePostgresRepo
}

func NewPostgresService(repos PrivateMessagePostgresRepo) *PostgresService {
	return &PostgresService{repos: repos}
}

func (s *PostgresService) SendPrivateMessage(receiver string, msg entities.Message) error {
	return s.repos.SendPrivateMessage(receiver, msg)
}

func (s *PostgresService) GetPrivateMessages(chat entities.ChatMetadata, params entities.PaginateParam) ([]entities.Message, error) {
	messages, err := s.repos.GetPrivateChat(chat)
	if err != nil {
		return nil, err
	}

	return pagination.Paginate(messages, params.Offset, params.Limit), nil
}

func (s *PostgresService) GetUsersWithMessage(receiver string) ([]string, error) {
	return s.repos.GetUserList(receiver)
}
