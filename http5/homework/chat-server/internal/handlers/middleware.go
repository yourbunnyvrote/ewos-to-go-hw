package handlers

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"
)

func (h *Handler) userIdentity(next http.Handler) http.Handler {
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
		if len(authCredentials) != 2 {
			http.Error(w, "Invalid Authorization credentials", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "username", authCredentials[0])
		ctx = context.WithValue(ctx, "password", authCredentials[1])
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) isUserExists(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value("username").(string)
		if !ok {
			http.Error(w, "Failed to get username from context", http.StatusInternalServerError)
			return
		}

		password, ok := r.Context().Value("password").(string)
		if !ok {
			http.Error(w, "Failed to get password from context", http.StatusInternalServerError)
			return
		}

		if _, err := h.service.GetUser(username, password); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
