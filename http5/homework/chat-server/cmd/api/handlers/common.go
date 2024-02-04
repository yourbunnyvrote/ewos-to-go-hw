package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

const MessagesPerPage = 3

func decodeRequestBody(r *http.Request, requestBody interface{}) error {
	return json.NewDecoder(r.Body).Decode(requestBody)
}

func sendResponse(w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func getPageParams(r *http.Request) (int, error) {
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	return page, err
}

func getPaginationIndexes(page int) (int, int) {
	startIndex := (page - 1) * MessagesPerPage
	endIndex := page * MessagesPerPage
	return startIndex, endIndex
}
