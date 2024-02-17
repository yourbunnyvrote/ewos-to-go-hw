package baseresponse

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	ErrorText string `json:"error,omitempty"`
}

func NewErrResponse(err error, statusCode int) *ErrResponse {
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

func RenderErr(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	respErr := NewErrResponse(err, statusCode)
	_ = render.Render(w, r, respErr)
}
