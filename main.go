package main

import (
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
	repo, err := db.New(utils.GetPostgresDSN())
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

	// Auth
	http.HandleFunc("/api/v1/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.RegisterHandler(w, r, repo, jwtGen, passwordHasher)
	})

	http.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.LoginHandler(w, r, repo, jwtGen, passwordHasher)
	})

	http.HandleFunc("/api/v1/offers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.GetOffersHandler(w, r, repo)
	})

	log.Printf("Starting server on port %s with DB user %s", port, dbUser)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
