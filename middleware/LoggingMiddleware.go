package middleware

import (
	"net/http"
	"passkey-server/utils/logger"
	"time"
)

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		logger.Infof(
			"[REQUEST] Method: %s | Path: %s | Duration: %v | Started at: %s | IP: %s | User-Agent: %s",
			r.Method,
			r.URL.Path,
			duration,
			start.String(),
			r.RemoteAddr,
			r.UserAgent(),
		)
	}
}
