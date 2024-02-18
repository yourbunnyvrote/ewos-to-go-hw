package mapper

import (
	"github.com/ew0s/ewos-to-go-hw/internal/api/resposne"
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
)

func MakeEntityAuthCredentials(login, password string) entities.AuthCredentials {
	return entities.AuthCredentials{
		Login:    login,
		Password: password,
	}
}

func MakeAuthResponse(credentials entities.AuthCredentials) resposne.AuthResponse {
	return resposne.AuthResponse{
		Username: credentials.Login,
	}
}
