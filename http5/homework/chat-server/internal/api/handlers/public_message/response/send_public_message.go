package response

import "github.com/ew0s/ewos-to-go-hw/internal/domain/entities"

type SendPublicMessagesResponse struct {
	Message entities.Message `json:"message"`
}
