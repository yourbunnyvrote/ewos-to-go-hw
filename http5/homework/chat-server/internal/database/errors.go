package database

import "errors"

var (
	ErrorMsgIsEmpty        = errors.New("message is empty")
	ErrorUserAlreadyExists = errors.New("user already exists")
	ErrorUserNotFound      = errors.New("user not found")
	ErrorDataError         = errors.New("incorrect data in memory db")
	ErrorUsernameEmpty     = errors.New("field 'username' is empty")
	ErrorPasswordEmpty     = errors.New("field 'password' is empty")
)
