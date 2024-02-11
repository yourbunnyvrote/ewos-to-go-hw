package handlers

import (
	"fmt"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/api/mapper"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/api/models/request"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/apiutils"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/constants"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/httputils/baseresponse"
	"github.com/go-chi/chi"
	"net/http"
)

type PrivateChatting interface {
	SendPrivateMessage(chat entities.Chat, msg entities.Message) error
	GetPrivateMessages(chat entities.Chat) ([]entities.Message, error)
	GetUsersWithMessage(receiver string) ([]string, error)
}

type PrivateChatHandler struct {
	service      PrivateChatting
	userIdentity *UserIdentity
}

func NewPrivateChatHandler(service PrivateChatting, userIdentity *UserIdentity) *PrivateChatHandler {
	return &PrivateChatHandler{
		service:      service,
		userIdentity: userIdentity,
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
	var textMessage request.TextMessage

	err := apiutils.DecodeRequestBody(r, &textMessage)
	if err != nil {
		baseresponse.RenderErr(w, r, fmt.Errorf("%w: %s", constants.ErrBadRequest, err))
		return
	}

	receiver := r.URL.Query().Get(constants.UsernameQueryParameter)

	sender, ok := r.Context().Value(constants.RouteContextUsernameValue).(string)
	if !ok {
		baseresponse.RenderErr(w, r, constants.ErrSomeServer)
		return
	}

	chat := mapper.MakeChat(receiver, sender)

	msg := mapper.MakeMessage(sender, textMessage.Content)

	err = h.service.SendPrivateMessage(chat, msg)
	if err != nil {
		baseresponse.RenderErr(w, r, fmt.Errorf("%w: %s", constants.ErrBadRequest, err))
		return
	}

	err = apiutils.SendResponse(w, http.StatusOK, msg)
	if err != nil {
		baseresponse.RenderErr(w, r, err)
		return
	}
}

// ShowPrivateMessages
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
	receiver := r.URL.Query().Get(constants.UsernameQueryParameter)

	sender, ok := r.Context().Value(constants.RouteContextUsernameValue).(string)
	if !ok {
		baseresponse.RenderErr(w, r, constants.ErrSomeServer)
		return
	}

	chat := mapper.MakeChat(receiver, sender)

	messages, err := h.service.GetPrivateMessages(chat)
	if err != nil {
		baseresponse.RenderErr(w, r, fmt.Errorf("%w: %s", constants.ErrNotFound, err))
		return
	}

	limit, offset, err := GetPaginateParameters(w, r)
	if err != nil {
		baseresponse.RenderErr(w, r, fmt.Errorf("%w: %s", constants.ErrBadRequest, err))
		return
	}

	pageMessages := PaginateMessages(messages, limit, offset)

	err = apiutils.SendResponse(w, http.StatusOK, pageMessages)
	if err != nil {
		baseresponse.RenderErr(w, r, err)
		return
	}
}

// ShowUsersWithMessages
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
	username, ok := r.Context().Value(constants.RouteContextUsernameValue).(string)
	if !ok {
		baseresponse.RenderErr(w, r, constants.ErrSomeServer)
		return
	}

	usersList, err := h.service.GetUsersWithMessage(username)
	if err != nil {
		baseresponse.RenderErr(w, r, err)
		return
	}

	err = apiutils.SendResponse(w, http.StatusOK, usersList)
	if err != nil {
		baseresponse.RenderErr(w, r, err)
		return
	}
}
