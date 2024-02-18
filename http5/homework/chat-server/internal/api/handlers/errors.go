package handlers

import "errors"

var (
	ErrUnauthorized          = errors.New("unauthorized")
	ErrRetrievingDataContext = errors.New("error retrieving data from context")
	ErrEmptyReceiver         = errors.New("empty receiver username")
)
