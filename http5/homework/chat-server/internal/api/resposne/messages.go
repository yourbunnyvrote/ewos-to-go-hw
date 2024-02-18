package resposne

import "github.com/ew0s/ewos-to-go-hw/internal/domain/entities"

type MessagesResponse struct {
	Messages []entities.Message `json:"messages"`
}
