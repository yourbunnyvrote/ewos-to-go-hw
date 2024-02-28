package request

import "github.com/go-playground/validator/v10"

type RegistrationRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (u RegistrationRequest) Validate(validate *validator.Validate) error {
	return validate.Struct(u)
}
