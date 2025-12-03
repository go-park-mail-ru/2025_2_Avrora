package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	// Total number of gRPC requests
	grpcRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "service", "status"},
	)

	// Duration of gRPC requests in seconds
	grpcRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "gRPC request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "service", "status"},
	)

	// Total number of gRPC errors
	grpcErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_errors_total",
			Help: "Total number of gRPC errors",
		},
		[]string{"method", "service", "status"},
	)
)

// GRPCMetricsInterceptor is a gRPC server interceptor to collect metrics
func GRPCMetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// Call the actual gRPC handler
	resp, err := handler(ctx, req)

	// Calculate duration
	duration := time.Since(start).Seconds()

	// Extract method and service names from the full method string
	fullMethod := info.FullMethod // Format: "/package.service/Method"
	service := extractServiceName(fullMethod)
	method := extractMethodName(fullMethod)

	// Determine status code
	statusCode := status.Code(err).String()

	// Update metrics
	grpcRequestsTotal.WithLabelValues(method, service, statusCode).Inc()
	grpcRequestDuration.WithLabelValues(method, service, statusCode).Observe(duration)
	if err != nil {
		grpcErrorsTotal.WithLabelValues(method, service, statusCode).Inc()
	}

	return resp, err
}

// Helper function to extract service name from full method string
func extractServiceName(fullMethod string) string {
	parts := strings.Split(fullMethod, "/")
	if len(parts) >= 3 {
		return parts[2] // Service name is the third part
	}
	return "unknown_service"
}

// Helper function to extract method name from full method string
func extractMethodName(fullMethod string) string {
	parts := strings.Split(fullMethod, "/")
	if len(parts) >= 4 {
		return parts[3] // Method name is the fourth part
	}
	return "unknown_method"
}