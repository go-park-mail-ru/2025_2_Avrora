package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"go.uber.org/zap"
)

// handler, проверяющий успешный доступ
func okHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(userID))
}

func TestAuthMiddleware_Success(t *testing.T) {
	logger := log.New(zap.NewNop())
	jwtGen := utils.NewJwtGenerator("secret123")

	token, err := jwtGen.GenerateJWT("user123")
	if err != nil {
		t.Fatalf("ошибка генерации токена: %v", err)
	}

	middleware := AuthMiddleware(logger, jwtGen)
	handler := middleware(http.HandlerFunc(okHandler))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("ожидался статус 200, получен %d", rec.Code)
	}
	if rec.Body.String() != "user123" {
		t.Fatalf("ожидался userID=user123, получен %s", rec.Body.String())
	}
}

func TestAuthMiddleware_NoAuthHeader(t *testing.T) {
	logger := log.New(zap.NewNop())
	jwtGen := utils.NewJwtGenerator("secret123")

	middleware := AuthMiddleware(logger, jwtGen)
	handler := middleware(http.HandlerFunc(okHandler))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("ожидался 401, получен %d", rec.Code)
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	logger := log.New(zap.NewNop())
	jwtGen := utils.NewJwtGenerator("secret123")

	middleware := AuthMiddleware(logger, jwtGen)
	handler := middleware(http.HandlerFunc(okHandler))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Token something")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("ожидался 401, получен %d", rec.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	logger := log.New(zap.NewNop())
	jwtGen := utils.NewJwtGenerator("secret123")

	middleware := AuthMiddleware(logger, jwtGen)
	handler := middleware(http.HandlerFunc(okHandler))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer bad.token.value")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("ожидался 401, получен %d", rec.Code)
	}
}
