package baseresponse

import (
	"encoding/json"
	"net/http"
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
		return ErrorMarshalResponse
	}

	_, err = w.Write(resp)
	if err != nil {
		return ErrorWriteResponse
	}

	return nil
}
