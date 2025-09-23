package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2025_2_Avrora/db"
	"github.com/go-park-mail-ru/2025_2_Avrora/handlers"
	"github.com/go-park-mail-ru/2025_2_Avrora/utils"
)

func main() {
	utils.LoadEnv()
	port := os.Getenv("PORT")
	dbUser := os.Getenv("DB_USER")
	repo := db.NewRepo()
	if err := repo.Init(fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", dbUser, os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))); err != nil {
		log.Fatal("ошибка инициализации бд:", err)
	}
	utils.InitJWTKey()

	// Auth
	http.HandleFunc("/api/v1/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterHandler(w, r, repo)
	})
	http.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(w, r, repo)
	})
	http.HandleFunc("/api/v1/offers", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetOffersHandler(w, r, repo)
	})

	log.Printf("Starting server on port %s with DB user %s", port, dbUser)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}