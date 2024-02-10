package baseresponse

import (
	"errors"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/constants"
	"github.com/go-chi/render"
	"net/http"
)

type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	ErrorText string `json:"error,omitempty"`
}

func NewErrResponse(err error, statusCode int) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: statusCode,
		ErrorText:      err.Error(),
	}
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrBadRequest(err error) render.Renderer {
	return NewErrResponse(err, http.StatusBadRequest)
}

func ErrUnauthorized(err error) render.Renderer {
	return NewErrResponse(err, http.StatusUnauthorized)
}

func ErrNotFound(err error) render.Renderer {
	return NewErrResponse(err, http.StatusNotFound)
}

func ErrUnknown(err error) render.Renderer {
	return NewErrResponse(err, http.StatusInternalServerError)
}

func RenderErr(w http.ResponseWriter, r *http.Request, err error) {
	var respErr render.Renderer

	switch {
	case errors.Is(err, constants.ErrBadRequest):
		respErr = ErrBadRequest(err)

	case errors.Is(err, constants.ErrUnauthorized):
		respErr = ErrUnauthorized(err)

	case errors.Is(err, constants.ErrNotFound):
		respErr = ErrNotFound(err)

	default:
		respErr = ErrUnknown(err)
	}

	_ = render.Render(w, r, respErr)
}
