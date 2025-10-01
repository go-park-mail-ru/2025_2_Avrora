package main

import (
	"fmt"
	"github.com/go-park-mail-ru/2025_2_Avrora/db"
	"github.com/go-park-mail-ru/2025_2_Avrora/handlers"
	"github.com/go-park-mail-ru/2025_2_Avrora/middleware"
	"github.com/go-park-mail-ru/2025_2_Avrora/utils"
	"log"
	"net/http"
	"os"
)

func main() {
	utils.LoadEnv()
	port := os.Getenv("SERVER_PORT")
	cors_origin := os.Getenv("CORS_ORIGIN")

	// Правильная строка подключения с паролем
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), // Добавляем пароль
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	log.Printf("Подключение к БД: %s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	repo, err := db.New(connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
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

	mux.HandleFunc("/api/v1/logout", func(w http.ResponseWriter, r *http.Request) {
		handlers.LogoutHandler(w, r, jwtGen)
	})

	mux.HandleFunc("/api/v1/offers", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetOffersHandler(w, r, repo)
	})

	mux.Handle("/api/v1/image/", http.StripPrefix("/api/v1/image/", http.FileServer(http.Dir("image/"))))

	handlerWithCORS := middleware.CorsMiddleware(mux, cors_origin)

	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handlerWithCORS))
}
