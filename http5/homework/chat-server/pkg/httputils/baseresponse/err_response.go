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

func (e *ErrResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func RenderErr(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	respErr := NewErrResponse(err, statusCode)

	if errRender := render.Render(w, r, respErr); errRender != nil {
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}
