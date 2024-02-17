package handlers

import "errors"

var (
	ErrorUnauthorized       = errors.New("unauthorized")
	ErrorPaginateParameters = errors.New("invalid paginate parameters")
	ErrorSomeServer         = errors.New("some problem")
	ErrorEmptyReceiver      = errors.New("empty receiver username")
)
