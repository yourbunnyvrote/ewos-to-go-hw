package entities

import "github.com/go-playground/validator/v10"

type PaginateParam struct {
	Limit  int `validator:"gt=0"`
	Offset int `validator:"gt=0"`
}

func (p PaginateParam) Validate() error {
	validate := validator.New()

	err := validate.Struct(p)
	if err != nil {
		return err
	}

	return nil
}
