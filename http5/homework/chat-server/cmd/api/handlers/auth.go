package handlers

import (
	"net/http"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/api/request"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"
)

type AuthService interface {
	CreateUser(user entities.User) (string, int, error)
	GetUser(user entities.User) (string, int, error)
}

type AuthHandler struct {
	serv AuthService
}

func NewAuthHandler(serv AuthService) *AuthHandler {
	return &AuthHandler{serv: serv}
}

// Registration
//
//	@Summary		Register a new user
//	@Description	Creates a new user account based on the provided user data
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		request.User	true	"User object for registration"
//	@Success		200		{object}	string			"User successfully registered"
//	@Failure		400		{string}	string			"Invalid request body"
//	@Failure		500		{string}	string			"Create user error"
//	@Failure		500		{string}	string			"JSON encoding error"
//	@Router			/register [post]
func (h *AuthHandler) Registration(w http.ResponseWriter, r *http.Request) {
	var newUserJSON request.User

	err := decodeRequestBody(r, &newUserJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newUser := entities.User(newUserJSON)

	response, statusCode, err := h.serv.CreateUser(newUser)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	err = sendResponse(w, statusCode, response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Authentication
//
//	@Summary		Authenticate a user
//	@Description	Authenticates an existing user based on the provided credentials
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		request.User	true	"User object for authentication"
//	@Success		200		{object}	string			"User successfully authenticated"
//	@Failure		400		{string}	string			"Invalid request body"
//	@Failure		500		{string}	string			"Authentication error"
//	@Failure		500		{string}	string			"JSON encoding error"
//	@Router			/auth [post]
func (h *AuthHandler) Authentication(w http.ResponseWriter, r *http.Request) {
	var existingUserJSON request.User

	err := decodeRequestBody(r, &existingUserJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingUser := entities.User(existingUserJSON)

	response, statusCode, err := h.serv.GetUser(existingUser)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	err = sendResponse(w, statusCode, response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
