package mapper

import "github.com/ew0s/ewos-to-go-hw/internal/api/handlers/request"

func MakeAuthCredentialsRequest(login, password string) request.AuthCredentials {
	return request.AuthCredentials{
		Login:    login,
		Password: password,
	}
}
