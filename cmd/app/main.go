package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/handlers"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/middleware"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/db"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/usecase"
)

func main() {
	utils.LoadEnv()
	cors_origin := os.Getenv("CORS_ORIGIN")
	port := os.Getenv("SERVER_PORT")
	dbConn, err := db.New(utils.GetPostgresDSN())
	if err != nil {
		log.Fatal(err)
	}
	hasher, err := utils.NewPasswordHasher(os.Getenv("PASSWORD_PEPPER"))
	if err != nil {
		log.Fatal(err)
	}
	jwtService := utils.NewJwtGenerator(os.Getenv("JWT_SECRET"))

	// Repository layer
	userRepo := db.NewUserRepository(dbConn.GetDB())
	offerRepo := db.NewOfferRepository(dbConn.GetDB())

	// Usecase layer
	authUC := usecase.NewAuthUsecase(userRepo, hasher, jwtService)
	offerUC := usecase.NewOfferUsecase(offerRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authUC)
	offerHandler := handlers.NewOfferHandler(offerUC)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/register", authHandler.Register)

	mux.HandleFunc("/api/v1/login", authHandler.Login)

	mux.HandleFunc("/api/v1/logout", authHandler.Logout)

	mux.HandleFunc("/api/v1/offers", offerHandler.GetOffers)

	mux.Handle("/api/v1/image/", http.StripPrefix("/api/v1/image/", http.FileServer(http.Dir("image/"))))

	handlerWithCORS := middleware.CorsMiddleware(mux, cors_origin)

	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handlerWithCORS))
}
