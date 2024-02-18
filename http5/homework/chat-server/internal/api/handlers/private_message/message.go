package private_message

import (
	"net/http"

	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers"
	handlersMapper "github.com/ew0s/ewos-to-go-hw/internal/api/handlers/mapper"
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/private_message/request"
	"github.com/ew0s/ewos-to-go-hw/internal/api/mapper"
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/pkg/httputils"
	"github.com/ew0s/ewos-to-go-hw/pkg/httputils/baseresponse"

	"github.com/go-chi/chi"
)

type PrivateMessageService interface {
	SendPrivateMessage(chat entities.ChatMetadata, msg entities.Message) error
	GetPrivateMessages(chat entities.ChatMetadata) ([]entities.Message, error)
	GetUsersWithMessage(receiver string) ([]string, error)
	PaginateMessages(messages []entities.Message, param entities.PaginateParam) []entities.Message
}

type Identity interface {
	Identify(next http.Handler) http.Handler
}

type PrivateChatHandler struct {
	service      PrivateMessageService
	userIdentity Identity
}

func NewPrivateChatHandler(service PrivateMessageService, userIdentity Identity) *PrivateChatHandler {
	return &PrivateChatHandler{
		service:      service,
		userIdentity: userIdentity,
	}
}

func (h *PrivateChatHandler) Routes() *chi.Mux {
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
//	@Param			text		body	request.SendPrivateMessageRequest	true	"Text message object for sending private message"
//	@Param			username	query	string				true	"Username of the message receiver"
//	@Security		BasicAuth
//	@Success		200	{object}	response.SendPrivateMessageResponse	"Message successfully sent"
//	@Failure		400	{string}	string				"Bad Request: Invalid request body or missing receiver"
//	@Failure		400	{string}	string				"Bad Request: Send private message error"
//	@Failure		401	{string}	string				"Unauthorized: Missing or invalid API key"
//	@Failure		500	{string}	string				"Failed to retrieve data from the query context"
//	@Failure		500	{string}	string				"JSON encoding error"
//	@Router			/v1/messages/private [post]
func (h *PrivateChatHandler) SendPrivateMessage(w http.ResponseWriter, r *http.Request) {
	var req request.SendPrivateMessageRequest

	if err := httputils.DecodeRequestBody(r, &req); err != nil {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, err)
		return
	}

	if err := req.Validate(); err != nil {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, err)
		return
	}

	receiver := r.URL.Query().Get(UsernameQueryParameter)

	sender, ok := r.Context().Value(RouteContextCredentials).(entities.AuthCredentials)
	if !ok {
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, handlers.ErrRetrievingDataContext)
		return
	}

	chat := mapper.MakeChatMetadata(receiver, sender.Login)
	if err := chat.Validate(); err != nil {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, err)
	}

	msg := mapper.MakeEntityMessage(sender.Login, req.Content)

	if err := h.service.SendPrivateMessage(chat, msg); err != nil {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, err)
		return
	}

	response := mapper.MakeSendPrivateMessageResponse(msg)

	if err := baseresponse.SendResponse(w, http.StatusOK, response); err != nil {
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, err)
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
//	@Success		200	{object}	response.ShowPrivateMessagesResponse	"List of private messages"
//	@Failure		400	{string}	string				"Missing receiver"
//	@Failure		400	{string}	string				"Invalid paginate parameters"
//	@Failure		401	{string}	string				"Unauthorized: Missing or invalid API key"
//	@Failure		404	{string}	string				"There is no dialog with such a person"
//	@Failure		500	{string}	string				"Failed to retrieve data from the query context"
//	@Failure		500	{string}	string				"JSON encoding error"
//	@Router			/v1/messages/private [get]
func (h *PrivateChatHandler) ShowPrivateMessages(w http.ResponseWriter, r *http.Request) {
	receiver := r.URL.Query().Get(UsernameQueryParameter)

	sender, ok := r.Context().Value(RouteContextCredentials).(entities.AuthCredentials)
	if !ok {
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, handlers.ErrRetrievingDataContext)
		return
	}

	req := mapper.MakeShowPrivateMessageRequest(receiver, sender.Login)

	if err := req.Validate(); err != nil {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, handlers.ErrEmptyReceiver)
	}

	chatMetadata := mapper.MakeChatMetadata(receiver, sender.Login)

	messages, err := h.service.GetPrivateMessages(chatMetadata)
	if err != nil {
		baseresponse.RenderErr(w, r, http.StatusNotFound, err)
		return
	}

	params, err := handlersMapper.GetPaginateParameters(r)
	if err != nil {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, err)
		return
	}

	pageMessages := h.service.PaginateMessages(messages, params)

	response := mapper.MakeShowPrivateMessageResponse(pageMessages)

	if err = baseresponse.SendResponse(w, http.StatusOK, response); err != nil {
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
//	@Success		200	{object}	response.ShowUsersWithMessagesResponse	"List of usernames with received private messages"
//	@Failure		401	{string}	string		"Unauthorized: Missing or invalid API key"
//	@Failure		500	{string}	string		"Failed to retrieve data from the query context"
//	@Failure		500	{string}	string		"JSON encoding error"
//	@Router			/v1/messages/private/users [get]
func (h *PrivateChatHandler) ShowUsersWithMessages(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(RouteContextCredentials).(entities.AuthCredentials)
	if !ok {
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, handlers.ErrRetrievingDataContext)
		return
	}

	req := mapper.MakeShowUsersWithMessagesRequest(user)

	usersList, err := h.service.GetUsersWithMessage(req.Username)
	if err != nil {
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, err)
		return
	}

	response := mapper.MakeUserListResponse(usersList)

	if err = baseresponse.SendResponse(w, http.StatusOK, response); err != nil {
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, err)
		return
	}
}
