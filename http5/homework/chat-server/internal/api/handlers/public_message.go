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

type PublicChatting interface {
	SendPublicMessage(msg entities.Message) error
	GetPublicMessages() ([]entities.Message, error)
}

type PublicChatHandler struct {
	service      PublicChatting
	userIdentity *UserIdentity
}

func NewPublicChatHandler(service PublicChatting, userIdentity *UserIdentity) *PublicChatHandler {
	return &PublicChatHandler{
		service:      service,
		userIdentity: userIdentity,
	}
}

func (h *PublicChatHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(h.userIdentity.Identify)
	r.Post("/", h.SendPublicMessage)
	r.Get("/", h.ShowPublicMessage)

	return r
}

// SendPublicMessage
//
//	@Summary		Send a public chat message
//	@Description	Sends a public chat message using the provided text content
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Param			text	body	request.TextMessage	true	"Text message object for sending public message"
//	@Security		BasicAuth
//	@Success		200	{object}	entities.Message	"Message successfully sent"
//	@Failure		400	{string}	string				"Bad Request: Invalid request body"
//	@Failure		401	{string}	string				"Unauthorized: Missing or invalid API key"
//	@Failure		500	{string}	string				"Bad Request: Send public message error"
//	@Failure		500	{string}	string				"JSON encoding error"
//	@Router			/messages/public [post]
func (h *PublicChatHandler) SendPublicMessage(w http.ResponseWriter, r *http.Request) {
	var textMessage request.TextMessage

	err := apiutils.DecodeRequestBody(r, &textMessage)
	if err != nil {
		baseresponse.RenderErr(w, r, err)
		return
	}

	username, ok := r.Context().Value(constants.RouteContextUsernameValue).(string)
	if !ok {
		baseresponse.RenderErr(w, r, constants.ErrSomeServer)
		return
	}

	msg := mapper.MakeMessage(username, textMessage.Content)

	err = h.service.SendPublicMessage(msg)
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

// ShowPublicMessage
//
//	@Summary		Get public chat messages
//	@Description	Retrieves public chat messages with pagination support
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Param			limit	query	integer	true	"Limit number for pagination (positive integer)"
//	@Param			offset	query	integer	true	"Start number for pagination (positive integer)"
//	@Security		BasicAuth
//	@Success		200	{object}	[]entities.Message	"List of public messages"
//	@Failure		400	{string}	string				"Bad Request: Invalid page number"
//	@Failure		401	{string}	string				"Unauthorized: Missing or invalid API key"
//	@Failure		500	{string}	string				"Bad Request: Get public messages error"
//	@Router			/messages/public [get]
func (h *PublicChatHandler) ShowPublicMessage(w http.ResponseWriter, r *http.Request) {
	messages, err := h.service.GetPublicMessages()
	if err != nil {
		baseresponse.RenderErr(w, r, err)
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
