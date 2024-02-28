package request

import "github.com/go-playground/validator/v10"

type SendPublicMessageRequest struct {
	Content string `json:"content" validate:"required"`
}

func (mr *SendPublicMessageRequest) Validate(validate *validator.Validate) error {
	return validate.Struct(mr)
}
