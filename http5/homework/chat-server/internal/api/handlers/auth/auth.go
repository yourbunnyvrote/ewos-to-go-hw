package auth

import (
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/auth/request"
	"net/http"

	"github.com/ew0s/ewos-to-go-hw/pkg/httputils"
	"github.com/ew0s/ewos-to-go-hw/pkg/httputils/baseresponse"

	"github.com/ew0s/ewos-to-go-hw/internal/api/mapper"
	"github.com/go-chi/chi"

	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
)

type AuthService interface {
	CreateUser(user entities.AuthCredentials) error
	GetUser(username string) (entities.AuthCredentials, error)
}

type Handler struct {
	service AuthService
}

func NewHandler(service AuthService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/sign-up", h.Registration)

	return r
}

// Registration
//
//	@Summary		Register a new user
//	@Description	Creates a new user account based on the provided user data
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		request.User	true	"User object for registration"
//	@Success		201		{object}	string			"User successfully registered"
//	@Failure		400		{string}	string			"Invalid request body"
//	@Failure		400		{string}	string			"Create user error"
//	@Failure		500		{string}	string			"JSON encoding error"
//	@Router			/v1/auth/sign-up [post]
func (h *Handler) Registration(w http.ResponseWriter, r *http.Request) {
	var req request.RegistrationRequest

	if err := httputils.DecodeRequestBody(r, &req); err != nil {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, err)
		return
	}

	if err := req.Validate(); err != nil {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, err)
		return
	}

	credentials := mapper.MakeEntityAuthCredentials(req.Username, req.Password)

	if err := h.service.CreateUser(credentials); err != nil {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, err)
		return
	}

	response := mapper.MakeAuthResponse(credentials)

	if err := baseresponse.SendResponse(w, http.StatusCreated, response); err != nil {
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, err)
		return
	}
}
