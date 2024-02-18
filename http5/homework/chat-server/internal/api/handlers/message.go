package handlers

import (
	"net/http"
	"strconv"

	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
)

func GetPaginateParameters(r *http.Request) (entities.PaginateParam, error) {
	limit, offset := getQueryPaginateParam(r)

	parameters, err := makePaginateParam(limit, offset)
	if err != nil {
		return entities.PaginateParam{}, err
	}

	err = checkPaginateParam(parameters)
	if err != nil {
		return entities.PaginateParam{}, err
	}

	return parameters, nil
}

func getQueryPaginateParam(r *http.Request) (limit, offset string) {
	limit = r.URL.Query().Get(LimitQueryParameter)
	offset = r.URL.Query().Get(OffsetQueryParameter)

	return limit, offset
}

func makePaginateParam(limit, offset string) (entities.PaginateParam, error) {
	var (
		parameters entities.PaginateParam
		err        error
	)

	parameters.Limit, err = strconv.Atoi(limit)
	if err != nil {
		return entities.PaginateParam{}, err
	}

	parameters.Offset, err = strconv.Atoi(offset)
	if err != nil {
		return entities.PaginateParam{}, err
	}

	return parameters, nil
}

func checkPaginateParam(param entities.PaginateParam) error {
	if param.Limit < 0 || param.Offset < 0 {
		return ErrorPaginateParameters
	}

	return nil
}
