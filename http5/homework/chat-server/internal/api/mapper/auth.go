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

func MakeResponseAuthCredentials(credentials entities.AuthCredentials) resposne.AuthCredentialsResponse {
	return resposne.AuthCredentialsResponse{
		Login: credentials.Login,
	}
}
