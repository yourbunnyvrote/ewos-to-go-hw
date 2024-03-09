package mapper

import (
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/auth/response"
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/request"
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
)

func MakeAuthCredentialsRequest(login, password string) request.AuthCredentials {
	return request.AuthCredentials{
		Login:    login,
		Password: password,
	}
}

func MakeEntityAuthCredentials(login, password string) entities.AuthCredentials {
	return entities.AuthCredentials{
		Login:    login,
		Password: password,
	}
}

func MakeRegistrationResponse(credentials entities.AuthCredentials) response.RegistrationResponse {
	return response.RegistrationResponse{
		Username: credentials.Login,
	}
}

func MakeJWTResponse(token string) response.JWTResponse {
	return response.JWTResponse{
		Token: token,
	}
}
