package repository

import "errors"

var (
	ErrorUserAlreadyExists = errors.New("user already exists")
	ErrorUserNotFound      = errors.New("user not found")
	ErrorDataError         = errors.New("incorrect data in memory db")
)
