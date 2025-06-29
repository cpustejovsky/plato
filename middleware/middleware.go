package middleware

import (
	"io"
	"log/slog"
	"net/http"
)

type MiddlewareFunc func(http.Handler) http.Handler

// TODO: Turn into a struct with a ServeHTTP method to allow the passing in of a logger?
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("start", "method", r.Method, "path", r.URL.Path)
		defer slog.Info("end", "method", r.Method, "path", r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

func DrainAndClose(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			_, _ = io.Copy(io.Discard, r.Body)
			_ = r.Body.Close()
		},
	)
}
