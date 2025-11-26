package main

import (
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/db"
	service "github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/grpc"
	fileserverpb "github.com/go-park-mail-ru/2025_2_Avrora/proto/fileserver"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/handlers"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/middleware"
	request_id "github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/middleware/request"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/utils"
	logger "github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/usecase"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	grpcLogger := appLogger.With(zap.String("layer", "grpc"))

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
	offerRepo := db.NewOfferRepository(dbConn.GetDB(), repoLogger)
	profileRepo := db.NewProfileRepository(dbConn.GetDB(), repoLogger)
	complexRepo := db.NewHousingComplexRepository(dbConn.GetDB(), repoLogger)

	// Usecases
	offerUC := usecase.NewOfferUsecase(offerRepo, usecaseLogger)
	profileUC := usecase.NewProfileUsecase(profileRepo, hasher, usecaseLogger)
	complexUC := usecase.NewHousingComplexUsecase(complexRepo, usecaseLogger)

	// Handlers
	offerHandler := handlers.NewOfferHandler(offerUC, httpLogger)
	profileHandler := handlers.NewProfileHandler(profileUC, httpLogger)
	complexHandler := handlers.NewComplexHandler(complexUC, httpLogger)

	// Auth middleware helper
	authMW := func(h http.HandlerFunc) http.HandlerFunc {
		return middleware.AuthMiddleware(appLogger, jwtService)(h).ServeHTTP
	}

	// GRPC Clients
	authClient, err := service.NewAuthClient(":50051", grpcLogger)
	if err != nil {
		log.Fatal("failed to create auth client", zap.Error(err))
	}

	// Create raw gRPC connection for fileserver
	fileServerConn, err := grpc.NewClient(":50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to create file server connection", zap.Error(err))
	}

	fileServerClient := fileserverpb.NewFileServerClient(fileServerConn)

	mux := http.NewServeMux()

	// Auth handler
	authHandler := handlers.NewAuthHandler(authClient, httpLogger)

	// Image handler with the proper gRPC client
	imageHandler := handlers.NewImageHandler(fileServerClient, httpLogger, "http://localhost:8080")

	// ┌───────────────┐
	// │ Public routes │
	// └───────────────┘
	mux.HandleFunc("/api/v1/register", authHandler.Register)
	mux.HandleFunc("/api/v1/login", authHandler.Login)
	mux.HandleFunc("/api/v1/logout", authHandler.Logout)

	// ┌──────────────────┐
	// │ Protected routes │
	// └──────────────────┘

	// Image routes
	mux.HandleFunc("/api/v1/image/upload", authMW(imageHandler.UploadImage))
	mux.Handle("/api/v1/image/", imageHandler.ImageServer())

	// Offers
	mux.HandleFunc("/api/v1/offers", offerHandler.GetOffers)
	mux.HandleFunc("/api/v1/offers/create", authMW(offerHandler.CreateOffer))
	mux.HandleFunc("/api/v1/offers/", offerHandler.GetOffer)
	mux.HandleFunc("/api/v1/offers/delete/", authMW(offerHandler.DeleteOffer))
	mux.HandleFunc("/api/v1/offers/update/", authMW(offerHandler.UpdateOffer))
	mux.HandleFunc("/api/v1/offers/pricehistory/", offerHandler.GetOfferPriceHistory)
	mux.HandleFunc("/api/v1/offers/viewcount/", offerHandler.GetViewCount)
	mux.HandleFunc("/api/v1/offers/view/", offerHandler.ViewOffer)

	// Like tracking endpoints
	mux.HandleFunc("/api/v1/offers/like/", authMW(offerHandler.ToggleLike))
	mux.HandleFunc("/api/v1/offers/likecount/", offerHandler.GetLikeCount)
	mux.HandleFunc("/api/v1/offers/isliked/", authMW(offerHandler.IsOfferLiked))

	// Profile
	mux.HandleFunc("/api/v1/profile/", authMW(profileHandler.GetProfile))
	mux.HandleFunc("/api/v1/profile/update/", authMW(profileHandler.UpdateProfile))
	mux.HandleFunc("/api/v1/profile/security/", authMW(profileHandler.UpdateProfileSecurityByID))
	mux.HandleFunc("/api/v1/profile/email/", authMW(profileHandler.UpdateEmail))
	mux.HandleFunc("/api/v1/profile/myoffers/", authMW(offerHandler.GetMyOffers))

	// Complex
	mux.HandleFunc("/api/v1/complexes/list", complexHandler.ListComplexes)
	mux.HandleFunc("/api/v1/complexes/create", authMW(complexHandler.CreateComplex))
	mux.HandleFunc("/api/v1/complexes/", complexHandler.GetComplexByID)
	mux.HandleFunc("/api/v1/complexes/update/", authMW(complexHandler.UpdateComplex))
	mux.HandleFunc("/api/v1/complexes/delete/", authMW(complexHandler.DeleteComplex))

	// Middleware setup
	var handler http.Handler = mux
	handler = middleware.CorsMiddleware(handler, corsOrigin)
	handler = request_id.RequestIDMiddleware(handler)
	handler = middleware.LoggerMiddleware(appLogger)(handler)

	appLogger.Logger.Info("starting server", zap.String("port", port))
	appLogger.Logger.Fatal("server stopped", zap.Error(http.ListenAndServe(":"+port, handler)))
}