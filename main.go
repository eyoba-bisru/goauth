package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/eyoba-bisru/goauth/config"
	"github.com/eyoba-bisru/goauth/handlers"
	"github.com/eyoba-bisru/goauth/logging"
	"github.com/eyoba-bisru/goauth/middleware"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.InitOAuthConfig(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), os.Getenv("GOOGLE_REDIRECT_URI"))

	mux := http.NewServeMux()

	// Create application logger (writes to stdout + logs/server.log)
	appLogger, closer, err := logging.NewLogger("logs/server.log")
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer closer()

	wrappedMux := middleware.BetterLoggingMiddleware(appLogger, mux)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Alive"))
	})
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/callback", handlers.CallbackHandler)
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	mux.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8" />
		<title>goauth — Swagger UI</title>
		<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@4/swagger-ui.css" />
	</head>
	<body>
		<div id="swagger-ui"></div>
		<script src="https://unpkg.com/swagger-ui-dist@4/swagger-ui-bundle.js"></script>
		<script>
			window.onload = function() {
				SwaggerUIBundle({
					url: '/swagger.json',
					dom_id: '#swagger-ui',
				});
			};
		</script>
	</body>
</html>`)
	})

	mux.HandleFunc("/redoc", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8" />
		<title>goauth — ReDoc</title>
		<meta name="viewport" content="width=device-width, initial-scale=1">
	</head>
	<body>
		<redoc spec-url='/swagger.json'></redoc>
		<script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script>
	</body>
</html>`)
	})

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err) // port not available, permission denied, etc.
	}

	fmt.Println("Server is running on port 8080")

	log.Fatal(http.Serve(ln, wrappedMux))
}
