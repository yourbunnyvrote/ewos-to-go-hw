package auth

import (
	"net/http"

	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/auth/request"
	"github.com/ew0s/ewos-to-go-hw/internal/api/mapper"
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/pkg/httputils"
	"github.com/ew0s/ewos-to-go-hw/pkg/httputils/baseresponse"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

type AuthService interface {
	CreateUser(user entities.AuthCredentials) error
	GetUser(username string) (entities.AuthCredentials, error)
}

type AuthHandler struct {
	service  AuthService
	validate *validator.Validate
}

func NewAuthHandler(service AuthService, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{
		service:  service,
		validate: validate,
	}
}

func (h *AuthHandler) Routes() *chi.Mux {
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
//	@Param			user	body		request.RegistrationRequest	true	"User object for registration"
//	@Success		201		{object}	string			"User successfully registered"
//	@Failure		400		{string}	string			"Invalid request body"
//	@Failure		400		{string}	string			"Create user error"
//	@Failure		500		{string}	string			"JSON encoding error"
//	@Router			/v1/auth/sign-up [post]
func (h *AuthHandler) Registration(w http.ResponseWriter, r *http.Request) {
	var req request.RegistrationRequest

	if err := httputils.DecodeRequestBody(r, &req); err != nil {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, err)
		return
	}

	if err := req.Validate(h.validate); err != nil {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, err)
		return
	}

	credentials := mapper.MakeEntityAuthCredentials(req.Username, req.Password)

	if err := h.service.CreateUser(credentials); err != nil {
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, err)
		return
	}

	response := mapper.MakeAuthResponse(credentials)

	if err := baseresponse.SendResponse(w, http.StatusCreated, response); err != nil {
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, err)
		return
	}
}
