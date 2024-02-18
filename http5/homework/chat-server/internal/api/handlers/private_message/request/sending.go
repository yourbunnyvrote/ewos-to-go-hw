package request

import "github.com/go-playground/validator/v10"

type SendPrivateMessageRequest struct {
	Content string `json:"content" validate:"required"`
}

func (mr *SendPrivateMessageRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(mr)
	if err != nil {
		return err
	}

	return nil
}
