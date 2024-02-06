package handlers

import (
	_ "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/docs" // swagger doc
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service"
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	serv *service.Service
}

func NewHandler(serv *service.Service) *Handler {
	return &Handler{serv: serv}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
	))

	r.Route("/v1", func(r chi.Router) {
		r.Post("/register", h.Registration)
		r.Post("/auth", h.Authentication)

		r.Group(func(r chi.Router) {
			r.Use(userIdentity)
			r.Use(h.isUserExists)
			r.Post("/messages/public", h.SendPublicMessage)
			r.Get("/messages/public", h.ShowPublicMessage)
		})

		r.Group(func(r chi.Router) {
			r.Use(userIdentity)
			r.Use(h.isUserExists)
			r.Post("/messages/private", h.SendPrivateMessage)
			r.Get("/messages/private", h.ShowPrivateMessages)
			r.Get("/messages/users", h.ShowUsersWithMessages)
		})
	})

	return r
}
