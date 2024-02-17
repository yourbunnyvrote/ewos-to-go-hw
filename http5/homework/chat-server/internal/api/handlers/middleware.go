package handlers

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/ew0s/ewos-to-go-hw/pkg/httputils/baseresponse"

	"github.com/ew0s/ewos-to-go-hw/internal/api/mapper"
)

type UserIdentity struct {
	service AuthService
}

func NewUserIdentity(service AuthService) *UserIdentity {
	return &UserIdentity{
		service: service,
	}
}

func (h *UserIdentity) Identify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			baseresponse.RenderErr(w, r, http.StatusUnauthorized, ErrorUnauthorized)
			return
		}

		authParts := strings.Fields(authHeader)
		if len(authParts) != 2 || authParts[0] != "Basic" {
			baseresponse.RenderErr(w, r, http.StatusUnauthorized, ErrorUnauthorized)
			return
		}

		decodedAuth, err := base64.StdEncoding.DecodeString(authParts[1])
		if err != nil {
			baseresponse.RenderErr(w, r, http.StatusUnauthorized, ErrorUnauthorized)
			return
		}

		authCredentials := strings.Split(string(decodedAuth), ":")
		if len(authCredentials) != CountCredentials {
			baseresponse.RenderErr(w, r, http.StatusUnauthorized, ErrorUnauthorized)
			return
		}

		credentials := mapper.MakeAuthCredentials(authCredentials[0], authCredentials[1])

		err = h.service.Identify(credentials)
		if err != nil {
			baseresponse.RenderErr(w, r, http.StatusUnauthorized, err)
		}

		ctx := context.WithValue(r.Context(), RouteContextUsernameValue, credentials.Login)
		ctx = context.WithValue(ctx, RouteContextPasswordValue, credentials.Password)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
