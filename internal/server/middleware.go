package server

import (
	"log/slog"
	"net/http"
)

func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			clientIP := r.RemoteAddr

			logger.Info("Served request",
				"method", r.Method,
				"path", r.URL.Path,
				"client_ip", clientIP,
			)
		})
	}
}
