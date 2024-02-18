package request

import "github.com/go-playground/validator/v10"

type SendPrivateMessageRequest struct {
	Content string `json:"content" validate:"required"`
}

func (r *SendPrivateMessageRequest) Validate() error {
	validate := validator.New()

	if err := validate.Struct(r); err != nil {
		return err
	}

	return nil
}
