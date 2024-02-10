package handlers

import (
	"context"
	"encoding/base64"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/constants"
	"net/http"
	"strings"
)

func UserIdentity(next http.Handler) http.Handler {
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
		if len(authCredentials) != constants.CountCredentials {
			http.Error(w, "Invalid Authorization credentials", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), constants.RouteContextUsernameValue, authCredentials[0])
		ctx = context.WithValue(ctx, constants.RouteContextPasswordValue, authCredentials[1])
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
