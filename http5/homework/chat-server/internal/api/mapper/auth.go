package mapper

import "github.com/ew0s/ewos-to-go-hw/internal/domain/entities"

func MakeAuthCredentials(login, password string) entities.AuthCredentials {
	return entities.AuthCredentials{
		Login:    login,
		Password: password,
	}
}
