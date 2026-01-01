package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/eyoba-bisru/goauth/handlers"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/callback", handlers.CallbackHandler)

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err) // port not available, permission denied, etc.
	}

	fmt.Println("Server is running on port 8080")

	log.Fatal(http.Serve(ln, nil))
}
