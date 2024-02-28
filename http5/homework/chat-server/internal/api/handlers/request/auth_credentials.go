package request

import "github.com/go-playground/validator/v10"

type AuthCredentials struct {
	Login    string `validator:"required"`
	Password string `validator:"required"`
}

func (r AuthCredentials) Validate(validate *validator.Validate) error {
	return validate.Struct(r)
}
