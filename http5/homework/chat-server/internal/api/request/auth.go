package request

import "github.com/go-playground/validator/v10"

type User struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (u User) Validate() error {
	validate := validator.New()

	err := validate.Struct(u)
	if err != nil {
		return err
	}

	return nil
}
