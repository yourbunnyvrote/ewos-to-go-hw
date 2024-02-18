package request

import "github.com/go-playground/validator/v10"

type SendPublicMessageRequest struct {
	Content string `json:"content" validate:"required"`
}

func (mr *SendPublicMessageRequest) Validate() error {
	validate := validator.New()

	return validate.Struct(mr)
}
