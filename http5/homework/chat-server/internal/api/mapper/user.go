package mapper

import (
	"github.com/ew0s/ewos-to-go-hw/internal/api/resposne"
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
)

func MakeEntityUser(username, password string) entities.User {
	return entities.User{
		Username: username,
		Password: password,
	}
}

func MakeUserResponse(user entities.User) resposne.UserResponse {
	return resposne.UserResponse{
		Username: user.Username,
	}
}
