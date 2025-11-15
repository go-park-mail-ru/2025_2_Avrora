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
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func main() {
	utils.LoadEnv()

	log, err := zap.NewProduction()
	if err != nil {
		panic("failed to create logger")
	}
	defer log.Sync()

	appLogger := logger.New(log)

	httpLogger := appLogger.With(zap.String("layer", "http"))
	usecaseLogger := appLogger.With(zap.String("layer", "usecase"))
	repoLogger := appLogger.With(zap.String("layer", "repository"))

	corsOrigin := os.Getenv("CORS_ORIGIN")
	port := os.Getenv("SERVER_PORT")

	// Database
	dbConn, err := db.New(utils.GetPostgresDSN())
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}

	// Services
	hasher, err := utils.NewPasswordHasher(os.Getenv("PASSWORD_PEPPER"))
	if err != nil {
		log.Fatal("failed to create password hasher", zap.Error(err))
	}
	jwtService := utils.NewJwtGenerator(os.Getenv("JWT_SECRET"))

	// Repositories
	userRepo := db.NewUserRepository(dbConn.GetDB(), repoLogger)
	offerRepo := db.NewOfferRepository(dbConn.GetDB(), repoLogger)
	profileRepo := db.NewProfileRepository(dbConn.GetDB(), repoLogger)
	complexRepo := db.NewHousingComplexRepository(dbConn.GetDB(), repoLogger)
	supportTicketRepo := db.NewSupportTicketRepository(dbConn.GetDB(), repoLogger) // Add support ticket repository

	// Usecases
	authUC := usecase.NewAuthUsecase(userRepo, hasher, jwtService, usecaseLogger)
	offerUC := usecase.NewOfferUsecase(offerRepo, usecaseLogger)
	profileUC := usecase.NewProfileUsecase(profileRepo, hasher, usecaseLogger)
	complexUC := usecase.NewHousingComplexUsecase(complexRepo, usecaseLogger)
	supportTicketUC := usecase.NewSupportTicketUsecase(supportTicketRepo, usecaseLogger) // Add support ticket usecase

	// Handlers
	authHandler := handlers.NewAuthHandler(authUC, httpLogger)
	offerHandler := handlers.NewOfferHandler(offerUC, httpLogger)
	profileHandler := handlers.NewProfileHandler(profileUC, httpLogger)
	complexHandler := handlers.NewComplexHandler(complexUC, httpLogger)
	supportTicketHandler := handlers.NewSupportTicketHandler(supportTicketUC, httpLogger)

	// Auth middleware helper
	authMW := func(h http.HandlerFunc) http.HandlerFunc {
		return middleware.AuthMiddleware(appLogger, jwtService)(h).ServeHTTP
	}

	mux := http.NewServeMux()

	// ┌───────────────┐
	// │ Public routes │
	// └───────────────┘
	mux.HandleFunc("/api/v1/register", authHandler.Register)
	mux.HandleFunc("/api/v1/login", authHandler.Login)
	mux.HandleFunc("/api/v1/logout", authHandler.Logout)

	// ┌──────────────────┐
	// │ Protected routes │
	// └──────────────────┘

	//Offers
	mux.HandleFunc("/api/v1/offers", offerHandler.GetOffers)
	mux.HandleFunc("/api/v1/offers/create", authMW(offerHandler.CreateOffer))
	mux.HandleFunc("/api/v1/offers/", offerHandler.GetOffer)
	mux.HandleFunc("/api/v1/offers/delete/", authMW(offerHandler.DeleteOffer))
	mux.HandleFunc("/api/v1/offers/update/", authMW(offerHandler.UpdateOffer))

	//Profile
	mux.HandleFunc("/api/v1/profile/", authMW(profileHandler.GetProfile))
	mux.HandleFunc("/api/v1/profile/update/", authMW(profileHandler.UpdateProfile))
	mux.HandleFunc("/api/v1/profile/security/", authMW(profileHandler.UpdateProfileSecurityByID))
	mux.HandleFunc("/api/v1/profile/email/", authMW(profileHandler.UpdateEmail))
	mux.HandleFunc("/api/v1/profile/myoffers/", authMW(offerHandler.GetMyOffers))

	//Complex
	mux.HandleFunc("/api/v1/complexes/list", complexHandler.ListComplexes)
	mux.HandleFunc("/api/v1/complexes/create", authMW(complexHandler.CreateComplex))
	mux.HandleFunc("/api/v1/complexes/", complexHandler.GetComplexByID)
	mux.HandleFunc("/api/v1/complexes/update/", authMW(complexHandler.UpdateComplex))
	mux.HandleFunc("/api/v1/complexes/delete/", authMW(complexHandler.DeleteComplex))

	// Support Tickets
	mux.HandleFunc("/api/v1/support-tickets", supportTicketHandler.CreateSupportTicket)

	mux.HandleFunc("/api/v1/support-tickets/all/", authMW(supportTicketHandler.GetAllSupportTickets))
	mux.HandleFunc("/api/v1/support-tickets/my/", authMW(supportTicketHandler.GetUserSupportTickets))
	mux.HandleFunc("/api/v1/support-tickets/", authMW(supportTicketHandler.GetSupportTicketByID))
	mux.HandleFunc("/api/v1/support-tickets/delete/", authMW(supportTicketHandler.DeleteSupportTicket))

	mux.HandleFunc("/api/v1/admin/support-tickets", authMW(supportTicketHandler.ListAllSupportTickets))
	mux.HandleFunc("/api/v1/admin/support-tickets/status/", authMW(supportTicketHandler.UpdateSupportTicketStatus))

	// Protected image file server
	mux.Handle("/api/v1/image/", handlers.RestrictedImageServer("./image"))
	imageHandler := handlers.NewImageHandler(usecaseLogger, "http://localhost:8080", "./image")
	mux.HandleFunc("/api/v1/image/upload", authMW(imageHandler.UploadImage))

	var handler http.Handler = mux
	handler = middleware.CorsMiddleware(handler, corsOrigin)
	handler = request_id.RequestIDMiddleware(handler)
	handler = middleware.LoggerMiddleware(appLogger)(handler)

	appLogger.Logger.Info("starting server", zap.String("port", port))
	appLogger.Logger.Fatal("server stopped", zap.Error(http.ListenAndServe(":"+port, handler)))
}
