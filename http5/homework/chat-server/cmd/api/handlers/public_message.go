package handlers

import (
	"net/http"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/api/request"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"
)

type PublicChatting interface {
	SendPublicMessage(msg entities.Message) error
	GetPublicMessages() ([]entities.Message, error)
}

type PublicChatHandler struct {
	serv PublicChatting
}

func NewPublicChatHandler(serv PublicChatting) *PublicChatHandler {
	return &PublicChatHandler{serv: serv}
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

	err := decodeRequestBody(r, &textMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username, ok := r.Context().Value(RouteContextUsernameValue).(string)
	if !ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	msg := entities.Message{
		Username: username,
		Content:  textMessage.Content,
	}

	err = h.serv.SendPublicMessage(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = sendResponse(w, http.StatusOK, msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	messages, err := h.serv.GetPublicMessages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pageMessages, err := paginateMessages(r, messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = sendResponse(w, http.StatusOK, pageMessages)
	if err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
}
