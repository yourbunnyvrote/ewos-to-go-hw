package repository

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrDataError         = errors.New("incorrect data in memory db")
)
