package mapper

import (
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
)

func MakeUser(username, password string) entities.User {
	return entities.User{
		Username: username,
		Password: password,
	}
}
