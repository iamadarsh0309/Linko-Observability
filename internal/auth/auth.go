package auth

import (
	"context"
	"log/slog"
	"net/http"
)

type contextKey string

const userContextKey contextKey = "user"

var allowedUsers = map[string]string{
	"frodo":   "$2a$10$B6O/n6teuCzpuh66jrUAdeaJ3WvXcxRkzpN0x7H.di9G9e/NGb9Me",
	"samwise": "$2a$10$EWZpvYhUJtJcEMmm/IBOsOGIcpxUnGIVMRiDlN/nxl1RRwWGkJtty",
	// frodo: "ofTheNineFingers"
	// samwise: "theStrong"
	"saruman": "invalidFormat",
}

type PasswordValidator func(password, stored string) (bool, error)

func Middleware(logger *slog.Logger, validatePassword PasswordValidator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		stored, exists := allowedUsers[username]
		if !exists {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ok, err := validatePassword(password, stored)
		if err != nil {
			logger.Error("error validating password", "user", username, "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), userContextKey, username))
		next.ServeHTTP(w, r)
	})
}

func Username(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(userContextKey).(string)
	return username, ok && username != ""
}
