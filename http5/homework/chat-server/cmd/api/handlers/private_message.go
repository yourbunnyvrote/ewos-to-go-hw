package handlers

import (
	"net/http"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/api/request"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/database"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"
)

type PrivateChatting interface {
	SendPrivateMessage(chat database.Chat, msg entities.Message) error
	GetPrivateMessages(chat database.Chat) ([]entities.Message, error)
	GetUsersWithMessage(receiver string) ([]string, error)
}

type PrivateChatHandler struct {
	serv PrivateChatting
}

func NewPrivateChatHandler(serv PrivateChatting) *PrivateChatHandler {
	return &PrivateChatHandler{serv: serv}
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
//	@Failure		401	{string}	string				"Unauthorized: Missing or invalid API key"
//	@Failure		500	{string}	string				"Bad Request: Send private message error"
//	@Failure		500	{string}	string				"JSON encoding error"
//	@Router			/messages/private [post]
func (h *PrivateChatHandler) SendPrivateMessage(w http.ResponseWriter, r *http.Request) {
	var textMessage request.TextMessage

	err := decodeRequestBody(r, &textMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	receiver := r.URL.Query().Get(UsernameQueryParameter)

	sender, ok := r.Context().Value(RouteContextUsernameValue).(string)
	if !ok {
		http.Error(w, "Failed to get username from context", http.StatusInternalServerError)
		return
	}

	chat := database.Chat{
		User1: receiver,
		User2: sender,
	}

	msg := entities.Message{
		Username: sender,
		Content:  textMessage.Content,
	}

	err = h.serv.SendPrivateMessage(chat, msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = sendResponse(w, http.StatusOK, msg)
	if err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
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
//	@Failure		400	{string}	string				"Bad Request: Invalid page number or missing receiver"
//	@Failure		401	{string}	string				"Unauthorized: Missing or invalid API key"
//	@Failure		500	{string}	string				"Bad Request: Get private messages error"
//	@Router			/messages/private [get]
func (h *PrivateChatHandler) ShowPrivateMessages(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := getPageParams(r)
	if err != nil || limit <= 0 || offset <= 0 {
		http.Error(w, "Invalid page parameters", http.StatusBadRequest)
		return
	}

	receiver := r.URL.Query().Get(UsernameQueryParameter)

	sender, ok := r.Context().Value(RouteContextUsernameValue).(string)
	if !ok {
		http.Error(w, "Failed to get username from context", http.StatusInternalServerError)
		return
	}

	chat := database.Chat{
		User1: receiver,
		User2: sender,
	}

	messages, err := h.serv.GetPrivateMessages(chat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startIndex, endIndex := getPaginationIndexes(limit, offset)

	pageMessages, err := paginateMessages(messages, startIndex, endIndex)
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
//	@Failure		400	{string}	string		"Bad Request: Get users with messages error"
//	@Failure		500	{string}	string		"JSON encoding error"
//	@Router			/messages/users [get]
func (h *PrivateChatHandler) ShowUsersWithMessages(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(RouteContextUsernameValue).(string)
	if !ok {
		http.Error(w, "Failed to get username from context", http.StatusInternalServerError)
		return
	}

	usersList, err := h.serv.GetUsersWithMessage(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = sendResponse(w, http.StatusOK, usersList)
	if err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
}
