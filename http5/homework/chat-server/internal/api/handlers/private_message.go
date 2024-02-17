package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/ew0s/ewos-to-go-hw/internal/api/request"

	"github.com/ew0s/ewos-to-go-hw/pkg/httputils"
	"github.com/ew0s/ewos-to-go-hw/pkg/httputils/baseresponse"

	"github.com/ew0s/ewos-to-go-hw/internal/api/mapper"
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
	"github.com/go-chi/chi"
)

type PrivateChatting interface {
	SendPrivateMessage(chat entities.UsersPair, msg entities.Message) error
	GetPrivateMessages(chat entities.UsersPair) ([]entities.Message, error)
	GetUsersWithMessage(receiver string) ([]string, error)
	PaginateMessages(messages []entities.Message, param entities.PaginateParam) []entities.Message
}

type PrivateChatHandler struct {
	service      PrivateChatting
	userIdentity *UserIdentity
	validate     *validator.Validate
}

func NewPrivateChatHandler(service PrivateChatting, userIdentity *UserIdentity) *PrivateChatHandler {
	return &PrivateChatHandler{
		service:      service,
		userIdentity: userIdentity,
		validate:     validator.New(),
	}
}

func (h *PrivateChatHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(h.userIdentity.Identify)
	r.Post("/", h.SendPrivateMessage)
	r.Get("/", h.ShowPrivateMessages)
	r.Get("/users", h.ShowUsersWithMessages)

	return r
}

// SendPrivateMessage
//
//	@Summary		Send a private chat message
//	@Description	Sends a private chat message to a specific user
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Param			text		body	request.TextMessage	true	"Text message object for sending private message"
//	@Param			username	query	string				true	"Username of the message receiver"
//	@Security		BasicAuth
//	@Success		200	{object}	entities.Message	"Message successfully sent"
//	@Failure		400	{string}	string				"Bad Request: Invalid request body or missing receiver"
//	@Failure		400	{string}	string				"Bad Request: Send private message error"
//	@Failure		401	{string}	string				"Unauthorized: Missing or invalid API key"
//	@Failure		500	{string}	string				"Failed to retrieve data from the query context"
//	@Failure		500	{string}	string				"JSON encoding error"
//	@Router			/messages/private [post]
func (h *PrivateChatHandler) SendPrivateMessage(w http.ResponseWriter, r *http.Request) {
	var req request.TextMessage

	err := httputils.DecodeRequestBody(r, &req)
	if err != nil {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, err)
		return
	}

	err = h.ValidateMessage(req)
	if err != nil {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, err)
		return
	}

	receiver := r.URL.Query().Get(UsernameQueryParameter)
	if receiver == "" {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, ErrorEmptyReceiver)
	}

	sender, ok := r.Context().Value(RouteContextUsernameValue).(string)
	if !ok {
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, ErrorSomeServer)
		return
	}

	chat := mapper.MakeChat(receiver, sender)

	msg := mapper.MakeMessage(sender, req.Content)

	err = h.service.SendPrivateMessage(chat, msg)
	if err != nil {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, err)
		return
	}

	err = baseresponse.SendResponse(w, http.StatusOK, msg)
	if err != nil {
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, ErrorSomeServer)
		return
	}
}

// ShowPrivateMessages
//
//	@Summary		Get private chat messages
//	@Description	Retrieves private chat messages with pagination support
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Param			limit		query	integer	true	"Limit number for pagination (positive integer)"
//	@Param			offset		query	integer	true	"Start number for pagination (positive integer)"
//	@Param			username	query	string	true	"Username of the message receiver"
//	@Security		BasicAuth
//	@Success		200	{object}	[]entities.Message	"List of private messages"
//	@Failure		400	{string}	string				"Missing receiver"
//	@Failure		400	{string}	string				"Invalid paginate parameters"
//	@Failure		401	{string}	string				"Unauthorized: Missing or invalid API key"
//	@Failure		404	{string}	string				"There is no dialog with such a person"
//	@Failure		500	{string}	string				"Failed to retrieve data from the query context"
//	@Failure		500	{string}	string				"JSON encoding error"
//	@Router			/messages/private [get]
func (h *PrivateChatHandler) ShowPrivateMessages(w http.ResponseWriter, r *http.Request) {
	receiver := r.URL.Query().Get(UsernameQueryParameter)
	if receiver == "" {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, ErrorEmptyReceiver)
	}

	sender, ok := r.Context().Value(RouteContextUsernameValue).(string)
	if !ok {
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, ErrorSomeServer)
		return
	}

	chat := mapper.MakeChat(receiver, sender)

	messages, err := h.service.GetPrivateMessages(chat)
	if err != nil {
		baseresponse.RenderErr(w, r, http.StatusNotFound, err)
		return
	}

	params, err := getPaginateParameters(r)
	if err != nil {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, err)
		return
	}

	pageMessages := h.service.PaginateMessages(messages, params)

	err = baseresponse.SendResponse(w, http.StatusOK, pageMessages)
	if err != nil {
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, err)
		return
	}
}

// ShowUsersWithMessages
//
//	@Summary		Get users with received private messages
//	@Description	Retrieves a list of users who have sent private messages to the authenticated user
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Security		BasicAuth
//	@Success		200	{object}	[]string	"List of usernames with received private messages"
//	@Failure		401	{string}	string		"Unauthorized: Missing or invalid API key"
//	@Failure		500	{string}	string		"Failed to retrieve data from the query context"
//	@Failure		500	{string}	string		"JSON encoding error"
//	@Router			/messages/private/users [get]
func (h *PrivateChatHandler) ShowUsersWithMessages(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(RouteContextUsernameValue).(string)
	if !ok {
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, ErrorSomeServer)
		return
	}

	usersList, err := h.service.GetUsersWithMessage(username)
	if err != nil {
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, ErrorSomeServer)
		return
	}

	err = baseresponse.SendResponse(w, http.StatusOK, usersList)
	if err != nil {
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, ErrorSomeServer)
		return
	}
}

func (h *PrivateChatHandler) ValidateMessage(req request.TextMessage) error {
	err := h.validate.Struct(req)
	if err != nil {
		return err
	}

	return nil
}
