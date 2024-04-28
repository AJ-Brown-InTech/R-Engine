package routes

import (
	"Engine/internal/models"
	"context"
	"net/http"
	"github.com/google/uuid"
)

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodedSessionToken := r.Header.Get("Authorization")

		// Check if the Authorization header exists
		var session = models.Session{Token: encodedSessionToken}

		if session.Token == "" {
			http.Error(w, "Unauthorized to make request, please authenticate", http.StatusUnauthorized)
			return
		}

		// Validate the session
		expired, err := session.Validate()
		if err != nil {
			http.Error(w, "Error validating session", http.StatusInternalServerError)
			return
		}

		if expired {
			http.Error(w, "Session token has expired, clearing token...", http.StatusUnauthorized)
			session.Token = ""
			r.Header.Set("Authorization","Bearer "+ session.Token)
		}
		next.ServeHTTP(w, r)
	})
}

func RequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rd := models.RequestData{}
		rd.From = r.Header.Get("Request")
		rd.TraceId = r.Header.Get("Trace")
		if rd.TraceId == "" {
			rd.TraceId = uuid.New().String()
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "Request", rd)))
	})
}