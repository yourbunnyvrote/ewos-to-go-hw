package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers"
	"github.com/ew0s/ewos-to-go-hw/internal/api/mapper"
	"github.com/ew0s/ewos-to-go-hw/pkg/httputils/baseresponse"
)

type JWTIdentityService interface {
	Identify(interface{}) (string, error)
}

type JWTIdentity struct {
	service JWTIdentityService
}

func NewJWTIdentity(service JWTIdentityService) *JWTIdentity {
	return &JWTIdentity{
		service: service,
	}
}

func (i *JWTIdentity) Identify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			baseresponse.RenderErr(w, r, http.StatusUnauthorized, handlers.ErrUnauthorized)
			return
		}

		authParts := strings.Fields(authHeader)
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			baseresponse.RenderErr(w, r, http.StatusUnauthorized, handlers.ErrUnauthorized)
			return
		}

		login, err := i.service.Identify(authParts[1])
		if err != nil {
			baseresponse.RenderErr(w, r, http.StatusInternalServerError, err)
			return
		}

		credentials := mapper.MakeEntityAuthCredentials(login, "")

		ctx := context.WithValue(r.Context(), RouteContextCredentials, credentials)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
