package request

import "github.com/go-playground/validator/v10"

type ShowPrivateMessageRequest struct {
	Receiver string `validator:"required"`
	Sender   string `validator:"required"`
}

func (r ShowPrivateMessageRequest) Validate() error {
	validate := validator.New()

	return validate.Struct(r)
}
