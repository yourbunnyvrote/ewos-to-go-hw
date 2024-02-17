package repository

import "github.com/ew0s/ewos-to-go-hw/internal/domain/entities"

type (
	UsersData        map[string]entities.User
	PublicChatsData  []entities.Message
	PrivateChatsData map[entities.UsersPair][]entities.Message
)
