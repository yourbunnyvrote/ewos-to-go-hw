package apiutils

import (
	"encoding/json"
	"net/http"
)

func DecodeRequestBody(r *http.Request, requestBody interface{}) error {
	return json.NewDecoder(r.Body).Decode(requestBody)
}
