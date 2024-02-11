package apiutils

import (
	"encoding/json"
	"net/http"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/constants"
)

type SuccessResponse struct {
	Result interface{} `json:"result"`
}

func SendResponse(w http.ResponseWriter, statusCode int, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	answer := SuccessResponse{Result: response}

	resp, err := json.Marshal(answer)
	if err != nil {
		return constants.ErrMarshalResponse
	}

	_, err = w.Write(resp)
	if err != nil {
		return constants.ErrWriteResponse
	}

	return nil
}
