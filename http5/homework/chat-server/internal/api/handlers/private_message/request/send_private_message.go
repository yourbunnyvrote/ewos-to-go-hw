package request

import "github.com/go-playground/validator/v10"

type SendPrivateMessageRequest struct {
	Content  string `json:"content" validate:"required"`
	Receiver string `validate:"required"`
}

func (r *SendPrivateMessageRequest) Validate() error {
	validate := validator.New()

	return validate.Struct(r)
}
