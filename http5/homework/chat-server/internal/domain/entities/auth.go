package entities

type AuthCredentials struct {
	Login    string `validate:"required"`
	Password string `validate:"required"`
}
