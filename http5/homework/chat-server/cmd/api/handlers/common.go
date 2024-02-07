package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

const (
	LimitQueryParameter  = "limit"
	OffsetQueryParameter = "offset"
)

var ErrEndOfPages = errors.New("no messages on this page")

func decodeRequestBody(r *http.Request, requestBody interface{}) error {
	return json.NewDecoder(r.Body).Decode(requestBody)
}

func sendResponse(w http.ResponseWriter, status int, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(response)
}
