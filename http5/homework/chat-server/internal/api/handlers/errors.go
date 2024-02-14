package handlers

import "errors"

var (
	ErrorBadRequest         = errors.New("bad request")
	ErrorNotFound           = errors.New("not found")
	ErrorUnauthorized       = errors.New("unauthorized")
	ErrorPaginateParameters = errors.New("invalid paginate parameters")
	ErrorSomeServer         = errors.New("some problem")
)
