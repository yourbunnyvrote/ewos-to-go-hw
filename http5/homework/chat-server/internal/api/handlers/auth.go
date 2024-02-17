package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/ew0s/ewos-to-go-hw/internal/api/request"

	"github.com/ew0s/ewos-to-go-hw/pkg/httputils"
	"github.com/ew0s/ewos-to-go-hw/pkg/httputils/baseresponse"

	"github.com/ew0s/ewos-to-go-hw/internal/api/mapper"
	"github.com/go-chi/chi"

	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
)

type AuthService interface {
	CreateUser(user entities.User) (string, error)
	GetUser(username string) (entities.User, error)
	Identify(user entities.AuthCredentials) error
}

type AuthHandler struct {
	service  AuthService
	validate *validator.Validate
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *AuthHandler) Routes() chi.Router {
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
//	@Router			/auth/sign-up [post]
func (h *AuthHandler) Registration(w http.ResponseWriter, r *http.Request) {
	var req request.User

	err := httputils.DecodeRequestBody(r, &req)
	if err != nil {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, err)
		return
	}

	err = h.ValidateUser(req)
	if err != nil {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, err)
		return
	}

	user := mapper.MakeUser(req.Username, req.Password)

	response, err := h.service.CreateUser(user)
	if err != nil {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, err)
		return
	}

	err = baseresponse.SendResponse(w, http.StatusCreated, response)
	if err != nil {
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, err)
		return
	}
}

func (h *AuthHandler) ValidateUser(req request.User) error {
	err := h.validate.Struct(req)
	if err != nil {
		return err
	}

	return nil
}
