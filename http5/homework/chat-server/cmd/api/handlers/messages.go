package handlers

import (
	"net/http"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/chatutil"
)

// SendPublicMessage
//
// @Summary		Send a public chat message
// @Description	Sends a public chat message using the provided text content
// @Tags			messages
// @Accept			json
// @Produce		json
// @Param			text	body	chatutil.TextMessage	true	"Text message object for sending public message"
// @Security		BasicAuth
// @Success		200	{object}	chatutil.Message	"Message successfully sent"
// @Failure		400	{string}	string				"Bad Request: Invalid request body"
// @Failure		401	{string}	string				"Unauthorized: Missing or invalid API key"
// @Failure		500	{string}	string				"Bad Request: Send public message error"
// @Failure		500	{string}	string				"JSON encoding error"
// @Router			/messages/public [post]
func (h *Handler) SendPublicMessage(w http.ResponseWriter, r *http.Request) {
	var textMessage chatutil.TextMessage

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

	msg := chatutil.Message{
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

// SendPrivateMessage
//
//	@Summary		Send a private chat message
//	@Description	Sends a private chat message to a specific user
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Param			text		body	chatutil.TextMessage	true	"Text message object for sending private message"
//	@Param			receiver	query	string					true	"Username of the message receiver"
//	@Security		BasicAuth
//	@Success		200	{object}	chatutil.Message	"Message successfully sent"
//	@Failure		400	{string}	string				"Bad Request: Invalid request body or missing receiver"
//	@Failure		401	{string}	string				"Unauthorized: Missing or invalid API key"
//	@Failure		500	{string}	string				"Bad Request: Send private message error"
//	@Failure		500	{string}	string				"JSON encoding error"
//	@Router			/messages/private [post]
func (h *Handler) SendPrivateMessage(w http.ResponseWriter, r *http.Request) {
	var textMessage chatutil.TextMessage

	err := decodeRequestBody(r, &textMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	receiver := r.URL.Query().Get("receiver")

	sender, ok := r.Context().Value(RouteContextUsernameValue).(string)
	if !ok {
		http.Error(w, "Failed to get username from context", http.StatusInternalServerError)
		return
	}

	chat := chatutil.Chat{
		User1: receiver,
		User2: sender,
	}

	msg := chatutil.Message{
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

// ShowPublicMessage
//
//	@Summary		Get public chat messages
//	@Description	Retrieves public chat messages with pagination support
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Param			page	query	integer	true	"Page number for pagination (positive integer)"
//	@Security		BasicAuth
//	@Success		200	{object}	[]chatutil.Message	"List of public messages"
//	@Failure		400	{string}	string				"Bad Request: Invalid page number"
//	@Failure		401	{string}	string				"Unauthorized: Missing or invalid API key"
//	@Failure		500	{string}	string				"Bad Request: Get public messages error"
//	@Router			/messages/public [get]
func (h *Handler) ShowPublicMessage(w http.ResponseWriter, r *http.Request) {
	page, err := getPageParams(r)
	if err != nil || page <= 0 {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	messages, err := h.serv.GetPublicMessages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startIndex, endIndex := getPaginationIndexes(page)

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

// ShowPrivateMessages
//
//	@Summary		Get private chat messages
//	@Description	Retrieves private chat messages with pagination support
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Param			page		query	integer	true	"Page number for pagination (positive integer)"
//	@Param			receiver	query	string	true	"Username of the message receiver"
//	@Security		BasicAuth
//	@Success		200	{object}	[]chatutil.Message	"List of private messages"
//	@Failure		400	{string}	string				"Bad Request: Invalid page number or missing receiver"
//	@Failure		401	{string}	string				"Unauthorized: Missing or invalid API key"
//	@Failure		500	{string}	string				"Bad Request: Get private messages error"
//	@Router			/messages/private [get]
func (h *Handler) ShowPrivateMessages(w http.ResponseWriter, r *http.Request) {
	page, err := getPageParams(r)
	if err != nil || page <= 0 {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	receiver := r.URL.Query().Get("receiver")

	sender, ok := r.Context().Value(RouteContextUsernameValue).(string)
	if !ok {
		http.Error(w, "Failed to get username from context", http.StatusInternalServerError)
		return
	}

	chat := chatutil.Chat{
		User1: receiver,
		User2: sender,
	}

	messages, err := h.serv.GetPrivateMessages(chat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startIndex, endIndex := getPaginationIndexes(page)

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
func (h *Handler) ShowUsersWithMessages(w http.ResponseWriter, r *http.Request) {
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
