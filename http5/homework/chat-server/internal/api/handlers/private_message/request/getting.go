package request

import "github.com/go-playground/validator/v10"

type ShowPrivateMessageRequest struct {
	Receiver string `validator:"required"`
	Sender   string `validator:"required"`
}

func (r ShowPrivateMessageRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(r)
	if err != nil {
		return err
	}

	return nil
}
