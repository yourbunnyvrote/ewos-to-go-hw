package private_message

import (
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/internal/service/message"
)

type PrivateMessageRepo interface {
	SendPrivateMessage(chat entities.ChatMetadata, msg entities.Message) error
	GetPrivateChats() (map[entities.ChatMetadata][]entities.Message, error)
}

type Service struct {
	repos PrivateMessageRepo
}

func NewService(repos PrivateMessageRepo) *Service {
	return &Service{repos: repos}
}

func (cs *Service) SendPrivateMessage(chat entities.ChatMetadata, msg entities.Message) error {
	if chat.Username1 < chat.Username2 {
		chat.Username1, chat.Username2 = chat.Username2, chat.Username1
	}

	return cs.repos.SendPrivateMessage(chat, msg)
}

func (cs *Service) GetPrivateMessages(chat entities.ChatMetadata, params entities.PaginateParam) ([]entities.Message, error) {
	privateChatsData, err := cs.repos.GetPrivateChats()
	if err != nil {
		return nil, err
	}

	if chat.Username1 < chat.Username2 {
		chat.Username1, chat.Username2 = chat.Username2, chat.Username1
	}

	pageMessages := message.PaginateMessages(privateChatsData[chat], params)

	return pageMessages, nil
}

func (cs *Service) GetUsersWithMessage(receiver string) ([]string, error) {
	privateChatsData, err := cs.repos.GetPrivateChats()
	if err != nil {
		return nil, err
	}

	listUsers := make([]string, 0)

	for key := range privateChatsData {
		switch receiver {
		case key.Username1:
			listUsers = append(listUsers, key.Username2)
		case key.Username2:
			listUsers = append(listUsers, key.Username1)
		}
	}

	return listUsers, nil
}
