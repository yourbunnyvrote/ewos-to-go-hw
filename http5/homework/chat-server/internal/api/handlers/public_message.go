package handlers

import (
	"github.com/ew0s/ewos-to-go-hw/internal/api/request"
	"net/http"

	"github.com/ew0s/ewos-to-go-hw/pkg/httputils"
	"github.com/ew0s/ewos-to-go-hw/pkg/httputils/baseresponse"

	"github.com/ew0s/ewos-to-go-hw/internal/api/mapper"
	"github.com/go-chi/chi"

	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
)

type PublicChatting interface {
	SendPublicMessage(msg entities.Message) error
	GetPublicMessages() ([]entities.Message, error)
	PaginateMessages(messages []entities.Message, params entities.PaginateParam) []entities.Message
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
//	@Failure		400	{string}	string				"Bad Request: Send public message error"
//	@Failure		401	{string}	string				"Unauthorized: Missing or invalid API key"
//	@Failure		500	{string}	string				"JSON encoding error"
//	@Router			/messages/public [post]
func (h *PublicChatHandler) SendPublicMessage(w http.ResponseWriter, r *http.Request) {
	var textMessage request.TextMessage

	err := httputils.DecodeRequestBody(r, &textMessage)
	if err != nil {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, err)
		return
	}

	username, ok := r.Context().Value(RouteContextUsernameValue).(string)
	if !ok {
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, ErrorSomeServer)
		return
	}

	msg := mapper.MakeMessage(username, textMessage.Content)

	err = h.service.SendPublicMessage(msg)
	if err != nil {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, err)
		return
	}

	err = baseresponse.SendResponse(w, http.StatusOK, msg)
	if err != nil {
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, err)
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
//	@Failure		400	{string}	string				"Invalid paginate parameters"
//	@Failure		401	{string}	string				"Unauthorized: Missing or invalid API key"
//	@Failure		500	{string}	string				"Bad Request: Get public messages error"
//	@Router			/messages/public [get]
func (h *PublicChatHandler) ShowPublicMessage(w http.ResponseWriter, r *http.Request) {
	messages, err := h.service.GetPublicMessages()
	if err != nil {
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, ErrorSomeServer)
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
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, ErrorSomeServer)
		return
	}
}
