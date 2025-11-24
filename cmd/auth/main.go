package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/db"
	service "github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/grpc"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/utils"
	logger "github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/usecase"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	utils.LoadEnv()

	log, err := zap.NewProduction()
	if err != nil {
		panic("failed to create logger")
	}
	defer log.Sync()

	appLogger := logger.New(log)

	usecaseLogger := appLogger.With(zap.String("layer", "usecase"))
	repoLogger := appLogger.With(zap.String("layer", "repository"))
	grpcLogger := appLogger.With(zap.String("service", "auth"))

	// Database
	dbConn, err := db.New(utils.GetPostgresDSN())
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer dbConn.Close() // Ensure database connection is closed on shutdown
	
	// Services
	hasher, err := utils.NewPasswordHasher(os.Getenv("PASSWORD_PEPPER"))
	if err != nil {
		log.Fatal("failed to create password hasher", zap.Error(err))
	}
	jwtService := utils.NewJwtGenerator(os.Getenv("JWT_SECRET"))
	authRepo := db.NewUserRepository(dbConn.GetDB(), repoLogger)
	authUC := usecase.NewAuthUsecase(authRepo, hasher, jwtService, usecaseLogger)

	grpcServer := grpc.NewServer()
	service.RegisterAuthServer(grpcServer, authUC, grpcLogger)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("failed to listen for gRPC", zap.Error(err))
	}

	// Graceful shutdown setup
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	appLogger.Logger.Info("gRPC server starting", zap.String("port", ":50051"))

	// Start gRPC server in a goroutine
	go func() {
		if err := grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			appLogger.Logger.Error("gRPC server failed", zap.Error(err))
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()

	appLogger.Logger.Info("shutting down gRPC server...")

	// Create a context with timeout for graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Gracefully stop the gRPC server
	go func() {
		grpcServer.GracefulStop()
		appLogger.Logger.Info("gRPC server stopped gracefully")
	}()

	// Wait for either graceful shutdown completion or timeout
	select {
	case <-shutdownCtx.Done():
		if shutdownCtx.Err() == context.DeadlineExceeded {
			appLogger.Logger.Warn("graceful shutdown timed out, forcing exit")
		}
	case <-time.After(1 * time.Second):
		appLogger.Logger.Warn("graceful shutdown timed out, forcing exit")
		grpcServer.Stop()
	}

	// Close database connection
	dbConn.Close()

	appLogger.Logger.Info("auth service shutdown complete")
}