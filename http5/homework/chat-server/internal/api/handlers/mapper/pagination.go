package mapper

import (
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
	"net/http"
	"strconv"
)

func GetPaginateParameters(r *http.Request) (entities.PaginateParam, error) {
	limit := r.URL.Query().Get(LimitQueryParameter)
	offset := r.URL.Query().Get(OffsetQueryParameter)

	parameters, err := makePaginateParam(limit, offset)
	if err != nil {
		return entities.PaginateParam{}, err
	}

	if err = parameters.Validate(); err != nil {
		return entities.PaginateParam{}, err
	}

	return parameters, nil
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
