package mapper

import "github.com/ew0s/ewos-to-go-hw/internal/api/resposne"

func MakeUserListResponse(usernames []string) resposne.UserListResponse {
	return resposne.UserListResponse{
		Usernames: usernames,
	}
}
