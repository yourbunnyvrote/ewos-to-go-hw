package handlers

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service"
	"github.com/go-chi/chi"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/users/register", h.Registration)
	r.Post("/users/auth", h.Authentication)

	return r
}
