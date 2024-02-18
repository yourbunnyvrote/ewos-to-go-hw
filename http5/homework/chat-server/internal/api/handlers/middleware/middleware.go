package middleware

import (
	"context"
	"encoding/base64"
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers"
	"net/http"
	"strings"

	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
	"github.com/go-playground/validator/v10"

	"github.com/ew0s/ewos-to-go-hw/pkg/httputils/baseresponse"

	"github.com/ew0s/ewos-to-go-hw/internal/api/mapper"
)

type Identity interface {
	Identify(user entities.AuthCredentials) error
}

type UserIdentity struct {
	service  Identity
	validate *validator.Validate
}

func NewUserIdentity(service Identity) *UserIdentity {
	return &UserIdentity{
		service:  service,
		validate: validator.New(),
	}
}

func (h *UserIdentity) Identify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			baseresponse.RenderErr(w, r, http.StatusUnauthorized, handlers.ErrorUnauthorized)
			return
		}

		authParts := strings.Fields(authHeader)
		if len(authParts) != 2 || authParts[0] != "Basic" {
			baseresponse.RenderErr(w, r, http.StatusUnauthorized, handlers.ErrorUnauthorized)
			return
		}

		decodedAuth, err := base64.StdEncoding.DecodeString(authParts[1])
		if err != nil {
			baseresponse.RenderErr(w, r, http.StatusUnauthorized, handlers.ErrorUnauthorized)
			return
		}

		authCredentials := strings.Split(string(decodedAuth), ":")
		if len(authCredentials) != handlers.CountCredentials {
			baseresponse.RenderErr(w, r, http.StatusUnauthorized, handlers.ErrorUnauthorized)
			return
		}

		credentials := mapper.MakeEntityAuthCredentials(authCredentials[0], authCredentials[1])

		err = h.ValidateCredentials(credentials)
		if err != nil {
			baseresponse.RenderErr(w, r, http.StatusUnauthorized, err)
			return
		}

		err = h.service.Identify(credentials)
		if err != nil {
			baseresponse.RenderErr(w, r, http.StatusUnauthorized, err)
			return
		}

		ctx := context.WithValue(r.Context(), handlers.RouteContextUsernameValue, credentials.Login)
		ctx = context.WithValue(ctx, handlers.RouteContextPasswordValue, credentials.Password)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (h *UserIdentity) ValidateCredentials(credentials entities.AuthCredentials) error {
	err := h.validate.Struct(credentials)
	if err != nil {
		return err
	}

	return nil
}
