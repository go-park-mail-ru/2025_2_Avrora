package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2025_2_Avrora/db"
	"github.com/go-park-mail-ru/2025_2_Avrora/handlers"
	"github.com/go-park-mail-ru/2025_2_Avrora/middleware"
	"github.com/go-park-mail-ru/2025_2_Avrora/utils"
)

func main() {
	utils.LoadEnv()
	port := os.Getenv("SERVER_PORT")
	cors_origin := os.Getenv("CORS_ORIGIN")
	dbUser := os.Getenv("DB_USER")
	repo, err := db.New(fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", dbUser, os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME")))
	if err != nil {
		log.Fatal(err)
	}
	jwtGen := utils.NewJwtGenerator(os.Getenv("JWT_SECRET"))
	pepper := os.Getenv("PASSWORD_PEPPER")
	if pepper == "" {
		log.Fatal("no pepper in .env")
	}

	passwordHasher, err := utils.NewPasswordHasher(pepper)
	if err != nil {
		log.Fatal("Ошибка инициализации хешера паролей:", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterHandler(w, r, repo, jwtGen, passwordHasher)
	})

	mux.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(w, r, repo, jwtGen, passwordHasher)
	})

	mux.HandleFunc("/api/v1/offers", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetOffersHandler(w, r, repo)
	})

	mux.Handle("/api/v1/image/", http.StripPrefix("/api/v1/image/", http.FileServer(http.Dir("image/"))))

	handlerWithCORS := middleware.CorsMiddleware(mux, cors_origin)

	log.Printf("Starting server on port %s with DB user %s", port, dbUser)
	log.Fatal(http.ListenAndServe(":"+port, handlerWithCORS))
}