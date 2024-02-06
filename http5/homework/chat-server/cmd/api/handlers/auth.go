package handlers

import (
	"net/http"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/chatutil"
)

// Registration
//
//	@Summary		Register a new user
//	@Description	Creates a new user account based on the provided user data
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		chatutil.User	true	"User object for registration"
//	@Success		200		{object}	string			"User successfully registered"
//	@Failure		400		{string}	string			"Invalid request body"
//	@Failure		500		{string}	string			"Create user error"
//	@Failure		500		{string}	string			"JSON encoding error"
//	@Router			/register [post]
func (h *Handler) Registration(w http.ResponseWriter, r *http.Request) {
	var newUser chatutil.User

	err := decodeRequestBody(r, &newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
//	@Param			user	body		chatutil.User	true	"User object for authentication"
//	@Success		200		{object}	string			"User successfully authenticated"
//	@Failure		400		{string}	string			"Invalid request body"
//	@Failure		500		{string}	string			"Authentication error"
//	@Failure		500		{string}	string			"JSON encoding error"
//	@Router			/auth [post]
func (h *Handler) Authentication(w http.ResponseWriter, r *http.Request) {
	var existingUser chatutil.User

	err := decodeRequestBody(r, &existingUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
