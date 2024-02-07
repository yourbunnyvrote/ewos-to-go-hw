package handlers

import (
	_ "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/docs" // swagger doc
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service"
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	auth *AuthHandler
	chat *ChattingHandler
}

func NewHandler(serv *service.Service) *Handler {
	return &Handler{
		auth: NewAuthHandler(serv.Auth),
		chat: NewChattingHandler(serv.Chat, serv.Chat),
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
	))

	r.Route("/v1", func(r chi.Router) {
		r.Post("/register", h.auth.Registration)
		r.Post("/auth", h.auth.Authentication)

		r.Group(func(r chi.Router) {
			r.Use(userIdentity)
			r.Use(h.isUserExists)
			r.Post("/messages/public", h.chat.public.SendPublicMessage)
			r.Get("/messages/public", h.chat.public.ShowPublicMessage)
		})

		r.Group(func(r chi.Router) {
			r.Use(userIdentity)
			r.Use(h.isUserExists)
			r.Post("/messages/private", h.chat.private.SendPrivateMessage)
			r.Get("/messages/private", h.chat.private.ShowPrivateMessages)
			r.Get("/messages/users", h.chat.private.ShowUsersWithMessages)
		})
	})

	return r
}
