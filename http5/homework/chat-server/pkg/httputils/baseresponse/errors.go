package baseresponse

import "errors"

var (
	ErrorMarshalResponse = errors.New("can't marshal response")
	ErrorWriteResponse   = errors.New("can't write response")
)
