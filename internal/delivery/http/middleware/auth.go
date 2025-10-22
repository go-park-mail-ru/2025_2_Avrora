package middleware

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/response"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"go.uber.org/zap"
)

type contextKey string

const UserContextKey contextKey = "user_id"

func AuthMiddleware(logger *log.Logger, jwtGen *utils.JwtGenerator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Error(r.Context(), "authorization header is empty")
				response.HandleError(w, nil, http.StatusUnauthorized, "требуется авторизация")
				return
			}

			const bearerPrefix = "Bearer "
			if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
				logger.Error(r.Context(), "invalid authorization header format")
				response.HandleError(w, nil, http.StatusUnauthorized, "некорректный формат токена")
				return
			}

			tokenStr := authHeader[len(bearerPrefix):]
			userID, err := jwtGen.ValidateJWT(tokenStr)
			if err != nil {
				logger.Error(r.Context(), "invalid token", zap.Error(err))
				response.HandleError(w, err, http.StatusUnauthorized, "недействительный токен")
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserContextKey).(string)
	return userID, ok
}