package entities

import "github.com/go-playground/validator/v10"

type ChatMetadata struct {
	Username1 string `validator:"required"`
	Username2 string `validator:"required"`
}

func (c ChatMetadata) Validate() error {
	validate := validator.New()

	return validate.Struct(c)
}
