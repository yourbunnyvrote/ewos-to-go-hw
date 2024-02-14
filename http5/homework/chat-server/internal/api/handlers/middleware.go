package handlers

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/httputils/baseresponse"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/api/mapper"
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
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		authParts := strings.Fields(authHeader)
		if len(authParts) != 2 || authParts[0] != "Basic" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		decodedAuth, err := base64.StdEncoding.DecodeString(authParts[1])
		if err != nil {
			http.Error(w, "Failed to decode Authorization header", http.StatusUnauthorized)
			return
		}

		authCredentials := strings.Split(string(decodedAuth), ":")
		if len(authCredentials) != CountCredentials {
			http.Error(w, "Invalid Authorization credentials", http.StatusUnauthorized)
			return
		}

		username, password := authCredentials[0], authCredentials[1]
		user1 := mapper.MakeUser(username, password)

		err = h.service.Identify(user1)
		if err != nil {
			baseresponse.RenderErr(w, r, err)
		}

		ctx := context.WithValue(r.Context(), RouteContextUsernameValue, username)
		ctx = context.WithValue(ctx, RouteContextPasswordValue, password)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
