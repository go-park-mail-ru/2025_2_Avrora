package middleware

import (
	"net"
	"net/http"
	"strings"
	"time"

	request_id "github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/middleware/request"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"go.uber.org/zap"
)

func LoggerMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			requestID, ok := r.Context().Value(request_id.RequestIDKey).(string)
			if !ok {
				requestID = request_id.GenerateRequestID()
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
				logger.Error(r.Context(),"server error", fields...)
			case ww.statusCode >= 400:
				logger.Warn(r.Context(), "client error", fields...)
			default:
				logger.Info(r.Context(), "request completed", fields...)
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