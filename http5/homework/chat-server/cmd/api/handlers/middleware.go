package handlers

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entities"
)

const (
	RouteContextUsernameValue = "username"
	RouteContextPasswordValue = "password"
	CountCredentials          = 2
)

func userIdentity(next http.Handler) http.Handler {
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

		ctx := context.WithValue(r.Context(), RouteContextUsernameValue, authCredentials[0])
		ctx = context.WithValue(ctx, RouteContextPasswordValue, authCredentials[1])
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) isUserExists(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value(RouteContextUsernameValue).(string)
		if !ok {
			http.Error(w, "Failed to get username from context", http.StatusInternalServerError)
			return
		}

		password, ok := r.Context().Value(RouteContextPasswordValue).(string)
		if !ok {
			http.Error(w, "Failed to get password from context", http.StatusInternalServerError)
			return
		}

		user := entities.User{
			Username: username,
			Password: password,
		}

		if _, statusCode, err := h.auth.serv.GetUser(user); err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		next.ServeHTTP(w, r)
	})
}
