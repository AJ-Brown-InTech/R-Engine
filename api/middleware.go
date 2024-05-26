package api

import (
	"Engine/types"
	"context"
	"net/http"
	"strings"
	"github.com/google/uuid"
)

// Define a custom type for context keys to avoid collisions.
type contextKey string

// ContextKeyUser is the key used to store the username in the context.
const ContextKeyUser = contextKey("user")

// GetUsernameFromRequest retrieves the username from the context.
func GetUsernameFromRequest(r *http.Request) (string, bool) {
	username, ok := r.Context().Value(ContextKeyUser).(string)
	return username, ok
}

// SessionMiddleware handles session validation for incoming requests.
func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: no Authorization header", http.StatusUnauthorized)
			return
		}

		// Extract the session token
		encodedSessionToken := strings.TrimSpace(authHeader)
		if encodedSessionToken == "" {
			http.Error(w, "Unauthorized: empty session token", http.StatusUnauthorized)
			return
		}

		var session types.Session

		// Validate the session
		valid, err := session.Validate(encodedSessionToken)
		if err != nil {
			http.Error(w, "Error validating session: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if !valid {
			http.Error(w, "Session expired or invalid", http.StatusUnauthorized)
			return
		}

		// Store the username in the request context
		ctx := context.WithValue(r.Context(), ContextKeyUser, session.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestMiddleware handles adding request data to the context.
func RequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rd := types.RequestData{}
		rd.From = r.Header.Get("Request")
		rd.TraceId = r.Header.Get("Trace")
		if rd.TraceId == "" {
			rd.TraceId = uuid.New().String()
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKey("Request"), rd)))
	})
}
