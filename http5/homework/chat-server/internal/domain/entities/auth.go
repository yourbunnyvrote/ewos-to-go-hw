package entities

import "github.com/go-playground/validator/v10"

type AuthCredentials struct {
	Login    string `validate:"required"`
	Password string `validate:"required"`
}

func (ac AuthCredentials) Validate() error {
	validate := validator.New()

	err := validate.Struct(ac)
	if err != nil {
		return err
	}

	return nil
}
