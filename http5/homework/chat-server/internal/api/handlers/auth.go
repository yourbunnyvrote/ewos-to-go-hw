package handlers

import (
	"fmt"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/api/mapper"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/apiutils"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/constants"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/httputils/baseresponse"
	"github.com/go-chi/chi"
	"net/http"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/api/request"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"
)

type AuthService interface {
	CreateUser(user entities.User) (string, error)
	GetUser(username string) (entities.User, error)
}

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/sign-up", h.Registration)
	r.Post("/sign-in", h.Authentication)

	return r
}

// Registration
// @Summary		Register a new user
// @Description	Creates a new user account based on the provided user data
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			user	body		request.User	true	"User object for registration"
// @Success		200		{object}	string			"User successfully registered"
// @Failure		400		{string}	string			"Invalid request body"
// @Failure		500		{string}	string			"Create user error"
// @Failure		500		{string}	string			"JSON encoding error"
// @Router			/register [post]
func (h *AuthHandler) Registration(w http.ResponseWriter, r *http.Request) {
	var user request.User

	err := apiutils.DecodeRequestBody(r, &user)
	if err != nil {
		baseresponse.RenderErr(w, r, fmt.Errorf("%w: %s", constants.ErrBadRequest, err))
		return
	}

	newUser := mapper.MakeUser(user.Username, user.Password)

	response, err := h.service.CreateUser(newUser)
	if err != nil {
		baseresponse.RenderErr(w, r, fmt.Errorf("%w: %s", constants.ErrBadRequest, err))
		return
	}

	err = apiutils.SendResponse(w, http.StatusCreated, response)
	if err != nil {
		baseresponse.RenderErr(w, r, err)
		return
	}
}

// Authentication
// @Summary		Authenticate a user
// @Description	Authenticates an existing user based on the provided credentials
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			user	body		request.User	true	"User object for authentication"
// @Success		200		{object}	string			"User successfully authenticated"
// @Failure		400		{string}	string			"Invalid request body"
// @Failure		500		{string}	string			"Authentication error"
// @Failure		500		{string}	string			"JSON encoding error"
// @Router			/auth [post]
func (h *AuthHandler) Authentication(w http.ResponseWriter, r *http.Request) {
	var user request.User

	err := apiutils.DecodeRequestBody(r, &user)
	if err != nil {
		baseresponse.RenderErr(w, r, err)
		return
	}

	existingUser := mapper.MakeUser(user.Username, user.Password)

	response, err := h.service.GetUser(existingUser.Username)
	if err != nil {
		baseresponse.RenderErr(w, r, err)
		return
	}

	err = apiutils.SendResponse(w, http.StatusOK, response)
	if err != nil {
		baseresponse.RenderErr(w, r, err)
		return
	}
}
