package mapper

import (
	"net/http"
	"strconv"

	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
)

func MakePaginateParam(limit, offset int) entities.PaginateParam {
	return entities.PaginateParam{
		Limit:  limit,
		Offset: offset,
	}
}

func GetPaginateParameters(r *http.Request) (limit, offset int, err error) {
	limitStr := r.URL.Query().Get(LimitQueryParameter)
	offsetStr := r.URL.Query().Get(OffsetQueryParameter)

	limit, offset, err = atoiPaginateParam(limitStr, offsetStr)
	if err != nil {
		return limit, offset, err
	}

	return limit, offset, nil
}

func atoiPaginateParam(limitStr, offsetStr string) (limit, offset int, err error) {
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		return limit, offset, err
	}

	offset, err = strconv.Atoi(offsetStr)
	if err != nil {
		return 0, 0, err
	}

	return limit, offset, nil
}
