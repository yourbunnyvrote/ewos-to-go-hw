package repository

import "github.com/ew0s/ewos-to-go-hw/internal/domain/entities"

type (
	UsersData        map[string]entities.AuthCredentials
	PublicChatsData  []entities.Message
	PrivateChatsData map[entities.ChatMetadata][]entities.Message
)
