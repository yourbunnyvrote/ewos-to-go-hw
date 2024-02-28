package request

import "github.com/go-playground/validator/v10"

type ShowPrivateMessageRequest struct {
	Receiver string `validator:"required"`
	Sender   string `validator:"required"`
	Limit    int    `validator:"gt=0"`
	Offset   int    `validator:"gt=0"`
}

func (r ShowPrivateMessageRequest) Validate(validate *validator.Validate) error {
	return validate.Struct(r)
}
