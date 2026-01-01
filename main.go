package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var OauthConfig *oauth2.Config

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	url := OauthConfig.AuthCodeURL("random-state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := OauthConfig.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Token exchange failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("token:", token)

	client := OauthConfig.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	fmt.Println("resp:", resp)

	w.Write([]byte("Login successful"))
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	OauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/callback", CallbackHandler)

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err) // port not available, permission denied, etc.
	}

	fmt.Println("Server is running on port 8080")

	log.Fatal(http.Serve(ln, nil))
}
