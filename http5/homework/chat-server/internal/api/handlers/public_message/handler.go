package public_message

import (
	"net/http"

	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/middleware"

	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers"
	handlersMapper "github.com/ew0s/ewos-to-go-hw/internal/api/handlers/mapper"
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/private_message"
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/public_message/request"
	"github.com/ew0s/ewos-to-go-hw/internal/api/mapper"
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/pkg/httputils"
	"github.com/ew0s/ewos-to-go-hw/pkg/httputils/baseresponse"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

type PublicMessageService interface {
	SendPublicMessage(msg entities.Message) error
	GetPublicMessages(params entities.PaginateParam) ([]entities.Message, error)
}

type Identity interface {
	Identify(next http.Handler) http.Handler
}

type Handler struct {
	service      PublicMessageService
	userIdentity Identity
	validate     *validator.Validate
}

func NewHandler(service PublicMessageService, userIdentity Identity, validate *validator.Validate) *Handler {
	return &Handler{
		service:      service,
		userIdentity: userIdentity,
		validate:     validate,
	}
}

func (h *Handler) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger(), h.userIdentity.Identify)
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
//	@Param			text	body	request.SendPublicMessageRequest	true	"Text message object for sending public message"
//	@Security		BasicAuth
//	@Success		200	{object}	response.SendPublicMessagesResponse	"Message successfully sent"
//	@Failure		400	{string}	string				"Bad Request: Invalid request body"
//	@Failure		400	{string}	string				"Bad Request: Send public message error"
//	@Failure		401	{string}	string				"Unauthorized: Missing or invalid API key"
//	@Failure		500	{string}	string				"JSON encoding error"
//	@Router			/v1/messages/public [post]
func (h *Handler) SendPublicMessage(w http.ResponseWriter, r *http.Request) {
	var req request.SendPublicMessageRequest

	if err := httputils.DecodeRequestBody(r, &req); err != nil {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, err)
		return
	}

	if err := req.Validate(h.validate); err != nil {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, err)
		return
	}

	credentials, ok := r.Context().Value(private_message.RouteContextCredentials).(entities.AuthCredentials)
	if !ok {
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, handlers.ErrRetrievingDataContext)
		return
	}

	msg := mapper.MakeEntityMessage(credentials.Login, req.Content)

	if err := h.service.SendPublicMessage(msg); err != nil {
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, err)
		return
	}

	response := mapper.MakeSendPublicMessageResponse(msg)

	if err := baseresponse.SendResponse(w, http.StatusOK, response); err != nil {
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
//	@Success		200	{object}	response.ShowPublicMessageResponse	"List of public messages"
//	@Failure		400	{string}	string				"Invalid paginate parameters"
//	@Failure		401	{string}	string				"Unauthorized: Missing or invalid API key"
//	@Failure		500	{string}	string				"Bad Request: Get public messages error"
//	@Router			/v1/messages/public [get]
func (h *Handler) ShowPublicMessage(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := handlersMapper.GetPaginateParameters(r)
	if err != nil {
		baseresponse.RenderErr(w, r, http.StatusBadRequest, err)
		return
	}

	params := handlersMapper.MakePaginateParam(limit, offset)

	pageMessages, err := h.service.GetPublicMessages(params)
	if err != nil {
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, err)
		return
	}

	response := mapper.MakeShowPublicMessagesResponse(pageMessages)

	if err = baseresponse.SendResponse(w, http.StatusOK, response); err != nil {
		baseresponse.RenderErr(w, r, http.StatusInternalServerError, err)
		return
	}
}
