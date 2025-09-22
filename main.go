package main

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Avrora/db"
	"github.com/go-park-mail-ru/2025_2_Avrora/handlers"
)

func main() {
	db.InitDB("postgres://postgres:postgres@localhost/2025_2_Avrora?sslmode=disable")

	// Auth
	http.HandleFunc("/api/v1/register", handlers.RegisterHandler)
	http.HandleFunc("/api/v1/login", handlers.LoginHandler)
	http.HandleFunc("/api/v1/logout", handlers.LogoutHandler)

	// Offers
	http.HandleFunc("/api/v1/offers", handlers.GetOffersHandler)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}