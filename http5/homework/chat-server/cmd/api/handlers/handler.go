package handlers

import (
	_ "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/docs"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service"
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
	))

	r.Post("/register", h.Registration)
	r.Post("/auth", h.Authentication)

	r.Group(func(r chi.Router) {
		r.Use(h.userIdentity)
		r.Use(h.isUserExists)
		r.Post("/messages/public", h.SendPublicMessage)
		r.Get("/messages/public", h.GetPublicMessage)
	})

	r.Group(func(r chi.Router) {
		r.Use(h.userIdentity)
		r.Use(h.isUserExists)
		r.Post("/messages/private", h.SendPrivateMessage)
		r.Get("/messages/private", h.GetPrivateMessages)
		r.Get("/messages/users", h.GetUsersWithMessages)
	})

	return r
}
