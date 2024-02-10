package constants

import "errors"

var (
	ErrMsgIsEmpty        = errors.New("message is empty")
	ErrEndOfPages        = errors.New("no messages on this page")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrDataError         = errors.New("incorrect data in memory db")
	ErrMarshalResponse   = errors.New("can't marshal response")
	ErrWriteResponse     = errors.New("can't write response")
	ErrBadRequest        = errors.New("bad request")
	ErrNotFound          = errors.New("not found")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrEmptyCredentials  = errors.New("username or password is empty")
	ErrIncorrectPassword = errors.New("incorrect password")
)
