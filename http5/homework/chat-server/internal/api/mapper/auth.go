package mapper

import (
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/auth/response"
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/request"
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
)

func MakeEntityAuthCredentials(credentials request.AuthCredentials) entities.AuthCredentials {
	return entities.AuthCredentials{
		Login:    credentials.Login,
		Password: credentials.Password,
	}
}

func MakeAuthResponse(credentials entities.AuthCredentials) response.RegistrationResponse {
	return response.RegistrationResponse{
		Username: credentials.Login,
	}
}
