package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	service "github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/grpc"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/middleware"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/utils"
	logger "github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func startMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":8083", nil); err != nil {
			log.Fatal("failed to start metrics server", zap.Error(err))
		}
	}()
}

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

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.GRPCMetricsInterceptor), 
	)

	service.RegisterFileServerServer(grpcServer, grpcLogger)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		grpcLogger.Logger.Fatal("failed to listen", zap.Error(err), zap.String("port", port))
	}

	startMetricsServer()

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

	grpcLogger.Logger.Info("shutting down fileserver gRPC server...")
	grpcServer.GracefulStop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	select {
	case <-shutdownCtx.Done():
		grpcLogger.Logger.Info("fileserver gRPC server stopped")
	case <-time.After(1 * time.Second):
		grpcLogger.Logger.Warn("fileserver gRPC server forced shutdown")
	}
}