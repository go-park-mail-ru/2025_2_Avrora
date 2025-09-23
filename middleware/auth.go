package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2025_2_Avrora/response"
	"github.com/go-park-mail-ru/2025_2_Avrora/utils"
)

// AuthMiddleware проверяет JWT токен и передаёт userID в context.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSONError(w, http.StatusUnauthorized, "Нет токена авторизации")
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeJSONError(w, http.StatusUnauthorized, "не верная схема авторизации")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Валидируем токен
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, "невалидный токен")
			return
		}

		// Передаём userID в context
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next(w, r.WithContext(ctx))
	}
}

// writeJSONError — вспомогательная функция для отправки ошибки в формате JSON
func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response.NewErrorResp(message))
}