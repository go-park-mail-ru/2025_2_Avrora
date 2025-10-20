package middleware

import (
	"net"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

func LoggerMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			requestID, ok := r.Context().Value(RequestIDKey).(string)
			if !ok {
				requestID = generateRequestID()
			}

			ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(ww, r)

			duration := time.Since(start)

			ip := getClientIP(r)

			remoteAddr := r.RemoteAddr

			referer := r.Referer()

			fields := []zap.Field{
				zap.String("request_id", requestID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.statusCode),
				zap.Duration("duration", duration),
				zap.String("ip", ip),                   
				zap.String("remote_addr", remoteAddr),  
				zap.String("referer", referer),         
				zap.String("user_agent", r.UserAgent()),
			}

			switch {
			case ww.statusCode >= 500:
				logger.Error("server error", fields...)
			case ww.statusCode >= 400:
				logger.Warn("client error", fields...)
			default:
				logger.Info("request completed", fields...)
			}
		})
	}
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