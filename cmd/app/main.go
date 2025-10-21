package main

import (
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/db"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/handlers"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/middleware"
	request_id "github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/middleware/request"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/utils"
	logger "github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/usecase"
	"go.uber.org/zap"
)

func main() {
	utils.LoadEnv()
	var log *zap.Logger
	log, _ = zap.NewProduction()
	defer log.Sync()

	appLogger := logger.New(log)

	httpLogger := appLogger.With(zap.String("layer", "http"))
	usecaseLogger := appLogger.With(zap.String("layer", "usecase"))
	repoLogger := appLogger.With(zap.String("layer", "repository"))

	cors_origin := os.Getenv("CORS_ORIGIN")
	port := os.Getenv("SERVER_PORT")
	dbConn, err := db.New(utils.GetPostgresDSN())
	if err != nil {
		log.Fatal("failed to connect to database: ", zap.Error(err))
	}
	hasher, err := utils.NewPasswordHasher("some")
	if err != nil {
		log.Fatal("failed to create password hasher: ", zap.Error(err))
	}
	jwtService := utils.NewJwtGenerator(os.Getenv("JWT_SECRET"))

	// Repository layer
	userRepo := db.NewUserRepository(dbConn.GetDB(), repoLogger)
	offerRepo := db.NewOfferRepository(dbConn.GetDB(), repoLogger)

	// Usecase layer
	authUC := usecase.NewAuthUsecase(userRepo, hasher, jwtService, usecaseLogger)
	offerUC := usecase.NewOfferUsecase(offerRepo, usecaseLogger)

	// Handlers
	authHandler := handlers.NewAuthHandler(authUC, httpLogger)
	offerHandler := handlers.NewOfferHandler(offerUC, httpLogger)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/register", authHandler.Register)

	mux.HandleFunc("/api/v1/login", authHandler.Login)

	mux.HandleFunc("/api/v1/logout", authHandler.Logout)

	mux.HandleFunc("/api/v1/offers", offerHandler.GetOffers)

	mux.Handle("/api/v1/image/", http.StripPrefix("/api/v1/image/", http.FileServer(http.Dir("image/"))))

	var handler http.Handler = mux

	handler = middleware.CorsMiddleware(handler, cors_origin)
	handler = request_id.RequestIDMiddleware(handler)
	handler = middleware.LoggerMiddleware(appLogger)(handler)

	appLogger.Logger.Info("starting server", zap.String("port", port))
	appLogger.Logger.Fatal("server stopped", zap.Error(http.ListenAndServe(":"+port, handler)))
}