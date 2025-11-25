package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	service "github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/grpc"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/utils"
	logger "github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	utils.LoadEnv()

	// Setup logger
	log, err := zap.NewProduction()
	if err != nil {
		panic("failed to create logger")
	}
	defer log.Sync()

	appLogger := logger.New(log)
	grpcLogger := appLogger.With(zap.String("layer", "grpc"))

	// Get configuration from environment
	port := "50052"
	if port == "" {
		port = "50052"
	}

	storageDir := os.Getenv("FILESERVER_STORAGE_DIR")
	if storageDir == "" {
		storageDir = "./image"
	}

	baseURL := os.Getenv("FILESERVER_BASE_URL")
	if baseURL == "" {
		baseURL = "/api/v1/image"
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()
	service.RegisterFileServerServer(grpcServer, grpcLogger)

	// Start listening
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		grpcLogger.Logger.Fatal("failed to listen", zap.Error(err), zap.String("port", port))
	}

	// Graceful shutdown handling
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start server in goroutine
	go func() {
		grpcLogger.Logger.Info("fileserver gRPC server starting", 
			zap.String("port", port),
			zap.String("storage_dir", storageDir),
			zap.String("base_url", baseURL))
		
		if err := grpcServer.Serve(lis); err != nil {
			grpcLogger.Logger.Error("gRPC server failed", zap.Error(err))
			stop()
		}
	}()

	<-ctx.Done()

	// Graceful shutdown
	grpcLogger.Logger.Info("shutting down fileserver gRPC server...")
	grpcServer.GracefulStop()
	
	// Allow some time for cleanup
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	
	select {
	case <-shutdownCtx.Done():
		grpcLogger.Logger.Info("fileserver gRPC server stopped")
	case <-time.After(1 * time.Second):
		grpcLogger.Logger.Warn("fileserver gRPC server forced shutdown")
	}
}