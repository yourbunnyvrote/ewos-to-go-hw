package handlers

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/chatutil"
	"net/http"
)

func (h *Handler) SendPublicMessage(w http.ResponseWriter, r *http.Request) {
	var textMessage chatutil.TextMessage
	err := decodeRequestBody(r, &textMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg := chatutil.Message{
		Username: r.Context().Value("username").(string),
		Content:  textMessage,
	}

	err = h.service.SendPublicMessage(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = sendResponse(w, msg)
	if err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) SendPrivateMessage(w http.ResponseWriter, r *http.Request) {
	var textMessage chatutil.TextMessage
	err := decodeRequestBody(r, &textMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	receiver := r.URL.Query().Get("receiver")
	sender := r.Context().Value("username").(string)
	chat := chatutil.Chat{
		User1: receiver,
		User2: sender,
	}

	msg := chatutil.Message{
		Username: sender,
		Content:  textMessage,
	}

	err = h.service.SendPrivateMessage(chat, msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = sendResponse(w, msg)
	if err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
}

func handleMessages(w http.ResponseWriter, messages []chatutil.Message, startIndex, endIndex int) {
	if startIndex >= len(messages) {
		http.Error(w, "No messages on this page", http.StatusNotFound)
		return
	}

	if endIndex > len(messages) {
		endIndex = len(messages)
	}

	pageMessages := messages[startIndex:endIndex]

	err := sendResponse(w, pageMessages)
	if err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetPublicMessage(w http.ResponseWriter, r *http.Request) {
	page, err := getPageParams(r)
	if err != nil || page <= 0 {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	startIndex, endIndex := getPaginationIndexes(page)

	messages, err := h.service.GetPublicMessages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handleMessages(w, messages, startIndex, endIndex)
}

func (h *Handler) GetPrivateMessages(w http.ResponseWriter, r *http.Request) {
	page, err := getPageParams(r)
	if err != nil || page <= 0 {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	startIndex, endIndex := getPaginationIndexes(page)

	receiver := r.URL.Query().Get("receiver")
	sender := r.Context().Value("username").(string)
	chat := chatutil.Chat{
		User1: receiver,
		User2: sender,
	}

	messages, err := h.service.GetPrivateMessages(chat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handleMessages(w, messages, startIndex, endIndex)
}

//func (h *Handler) GetPublicMessage(w http.ResponseWriter, r *http.Request) {
//	pageStr := r.URL.Query().Get("page")
//	page, err := strconv.Atoi(pageStr)
//	if err != nil || page <= 0 {
//		http.Error(w, "Invalid page number", http.StatusBadRequest)
//		return
//	}
//
//	perPage := 3
//	startIndex := (page - 1) * perPage
//	endIndex := page * perPage
//
//	messages, err := h.service.GetPublicMessages()
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	if startIndex >= len(messages) {
//		http.Error(w, "No messages on this page", http.StatusNotFound)
//		return
//	}
//
//	if endIndex > len(messages) {
//		endIndex = len(messages)
//	}
//
//	pageMessages := messages[startIndex:endIndex]
//
//	err = sendResponse(w, pageMessages)
//	if err != nil {
//		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
//		return
//	}
//}
//
//func (h *Handler) GetPrivateMessages(w http.ResponseWriter, r *http.Request) {
//	pageStr := r.URL.Query().Get("page")
//	page, err := strconv.Atoi(pageStr)
//	if err != nil || page <= 0 {
//		http.Error(w, "Invalid page number", http.StatusBadRequest)
//		return
//	}
//
//	perPage := 3
//	startIndex := (page - 1) * perPage
//	endIndex := page * perPage
//
//	receiver := r.URL.Query().Get("receiver")
//	sender := r.Context().Value("username").(string)
//	chat := chatutil.Chat{
//		User1: receiver,
//		User2: sender,
//	}
//
//	messages, err := h.service.GetPrivateMessages(chat)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	if startIndex >= len(messages) {
//		http.Error(w, "No messages on this page", http.StatusNotFound)
//		return
//	}
//
//	if endIndex > len(messages) {
//		endIndex = len(messages)
//	}
//
//	pageMessages := messages[startIndex:endIndex]
//
//	err = sendResponse(w, pageMessages)
//	if err != nil {
//		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
//		return
//	}
//}

func (h *Handler) GetUsersWithMessages(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)

	usersList, err := h.service.GetUsersWithMessage(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = sendResponse(w, usersList)
	if err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
}
