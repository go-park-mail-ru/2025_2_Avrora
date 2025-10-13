package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

const RequestIDKey contextKey = "request_id"

func LoggerMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = generateRequestID()
			}

			ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
			r = r.WithContext(ctx)

			ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(ww, r)

			duration := time.Since(start)
			fields := []zap.Field{
				zap.String("request_id", requestID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.statusCode),
				zap.Duration("duration", duration),
				zap.String("ip", getClientIP(r)),
				zap.String("user_agent", r.UserAgent()),
			}

			if ww.statusCode >= 500 {
				logger.Error("server error", fields...)
			} else if ww.statusCode >= 400 {
				logger.Warn("client error", fields...)
			} else {
				logger.Info("request completed", fields...)
			}
		})
	}
}

func generateRequestID() string {
	return time.Now().Format("20060102150405")[:8]
}

func getClientIP(r *http.Request) string {
	for _, h := range []string{"X-Forwarded-For", "X-Real-IP"} {
		if ip := r.Header.Get(h); ip != "" {
			if idx := strings.Index(ip, ","); idx != -1 {
				ip = ip[:idx]
			}
			return strings.TrimSpace(ip)
		}
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}