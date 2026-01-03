package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

// BetterLoggingMiddleware returns a middleware handler that logs requests using the
// provided logger. This allows logs to be written to a file or any other writer.
func BetterLoggingMiddleware(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Initialize our custom wrapper
		wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		logger.Printf(
			"STATUS: %d | METHOD: %s | PATH: %s | DURATION: %s",
			wrapped.status,
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}
