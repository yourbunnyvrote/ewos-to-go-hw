package handlers

import (
	"net/http"
	"strconv"
)

func getPaginateParameters(r *http.Request) (limit int, offset int, err error) {
	limitStr := r.URL.Query().Get(LimitQueryParameter)

	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		return 0, 0, ErrorPaginateParameters
	}

	offsetStr := r.URL.Query().Get(OffsetQueryParameter)

	offset, err = strconv.Atoi(offsetStr)
	if err != nil {
		return 0, 0, ErrorPaginateParameters
	}

	if limit < 0 || offset < 0 {
		return 0, 0, ErrorPaginateParameters
	}

	return limit, offset, nil
}
