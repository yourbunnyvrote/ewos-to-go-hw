package handlers

import "errors"

var (
	ErrUnauthorized  = errors.New("unauthorized")
	ErrorSomeServer  = errors.New("some problem")
	ErrEmptyReceiver = errors.New("empty receiver username")
)
