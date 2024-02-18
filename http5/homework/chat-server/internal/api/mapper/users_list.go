package mapper

import (
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/private_message/request"
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/private_message/response"
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
)

func MakeShowUsersWithMessagesRequest(credentials entities.AuthCredentials) request.ShowUsersWithMessagesRequest {
	return request.ShowUsersWithMessagesRequest{
		Username: credentials.Login,
	}
}

func MakeUserListResponse(usernames []string) response.ShowUsersWithMessagesResponse {
	return response.ShowUsersWithMessagesResponse{
		Usernames: usernames,
	}
}
