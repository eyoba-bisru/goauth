package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/eyoba-bisru/goauth/config"
	"github.com/eyoba-bisru/goauth/handlers"
	"github.com/joho/godotenv"
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

func BetterLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Initialize our custom wrapper
		wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		log.Printf(
			"STATUS: %d | METHOD: %s | PATH: %s | DURATION: %s",
			wrapped.status,
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.InitOAuthConfig(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), os.Getenv("GOOGLE_REDIRECT_URI"))

	mux := http.NewServeMux()

	wrappedMux := BetterLoggingMiddleware(mux)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Default().Println("/health")
		w.Write([]byte("Alive"))
	})
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/callback", handlers.CallbackHandler)

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err) // port not available, permission denied, etc.
	}

	fmt.Println("Server is running on port 8080")

	log.Fatal(http.Serve(ln, wrappedMux))
}
