package response

import "github.com/ew0s/ewos-to-go-hw/internal/domain/entities"

type ShowPrivateMessagesResponse struct {
	Messages []entities.Message `json:"messages"`
}
