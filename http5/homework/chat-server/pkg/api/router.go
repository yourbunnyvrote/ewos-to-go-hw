package api

import (
	"github.com/go-chi/chi"
)

func MakeRoutes(basicPath string, routers map[string]chi.Router) chi.Router {
	r := chi.NewRouter()

	r.Route(basicPath, func(r chi.Router) {
		for pattern, router := range routers {
			r.Mount(pattern, router)
		}
	})

	return r
}
