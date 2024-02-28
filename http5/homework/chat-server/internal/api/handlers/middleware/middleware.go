package middleware

import (
	"context"
	"encoding/base64"
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/mapper"
	mapper2 "github.com/ew0s/ewos-to-go-hw/internal/api/mapper"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"

	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers"
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/pkg/httputils/baseresponse"
)

type IdentityService interface {
	Identify(user entities.AuthCredentials) error
}

type UserIdentity struct {
	service  IdentityService
	validate *validator.Validate
}

func NewUserIdentity(service IdentityService) *UserIdentity {
	return &UserIdentity{
		service: service,
	}
}

func (h *UserIdentity) Identify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			baseresponse.RenderErr(w, r, http.StatusUnauthorized, handlers.ErrUnauthorized)
			return
		}

		authParts := strings.Fields(authHeader)
		if len(authParts) != 2 || authParts[0] != "Basic" {
			baseresponse.RenderErr(w, r, http.StatusUnauthorized, handlers.ErrUnauthorized)
			return
		}

		decodedAuth, err := base64.StdEncoding.DecodeString(authParts[1])
		if err != nil {
			baseresponse.RenderErr(w, r, http.StatusUnauthorized, handlers.ErrUnauthorized)
			return
		}

		authCredentials := strings.Split(string(decodedAuth), ":")
		if len(authCredentials) != CountCredentials {
			baseresponse.RenderErr(w, r, http.StatusUnauthorized, handlers.ErrUnauthorized)
			return
		}

		req := mapper.MakeAuthCredentialsRequest(authCredentials[0], authCredentials[1])

		if err = req.Validate(h.validate); err != nil {
			baseresponse.RenderErr(w, r, http.StatusUnauthorized, err)
			return
		}

		credentials := mapper2.MakeEntityAuthCredentials(req)

		if err = h.service.Identify(credentials); err != nil {
			baseresponse.RenderErr(w, r, http.StatusInternalServerError, err)
			return
		}

		ctx := context.WithValue(r.Context(), RouteContextCredentials, credentials)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
