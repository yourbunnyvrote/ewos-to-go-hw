package baseresponse

import "errors"

var (
	ErrorMarshalResponse = errors.New("can't marshal response")
	ErrorWriteResponse   = errors.New("can't write response")
	ErrorBadRequest      = errors.New("bad request")
	ErrorNotFound        = errors.New("not found")
	ErrorUnauthorized    = errors.New("unauthorized")
)
