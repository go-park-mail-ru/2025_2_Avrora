package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/go-park-mail-ru/2025_2_Avrora/response"
	"github.com/go-park-mail-ru/2025_2_Avrora/utils"
)

// Тестовый обработчик — проверяет, что userID передан
func dummyHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")
	if userID == nil {
		http.Error(w, "userID not found in context", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func TestMain(m *testing.M) {
	// Устанавливаем тестовый JWT ключ
	utils.SetJWTSecret([]byte("test_secret_32_chars_min_for_tests"))

	code := m.Run()
	os.Exit(code)
}

func TestAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	AuthMiddleware(dummyHandler)(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var resp response.ErrorResp
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	if resp.Error != "Нет токена авторизации" {
		t.Errorf("Expected error 'Missing Authorization header', got '%s'", resp.Error)
	}
}

func TestAuthMiddleware_WrongSchema(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Token 12345") // Неправильная схема
	w := httptest.NewRecorder()

	AuthMiddleware(dummyHandler)(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var resp response.ErrorResp
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	if resp.Error != "не верная схема авторизации" {
		t.Errorf("Expected error 'не верная схема авторизации', got '%s'", resp.Error)
	}
}

func TestAuthMiddleware_EmptyToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer ") // Токен отсутствует
	w := httptest.NewRecorder()

	AuthMiddleware(dummyHandler)(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var resp response.ErrorResp
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	if resp.Error != "невалидный токен" {
		t.Errorf("Expected error 'невалидный токен, got '%s'", resp.Error)
	}
}

func TestAuthMiddleware_InvalidTokenSignature(t *testing.T) {
	// Генерируем токен с другим ключом
	invalidKey := []byte("wrong_secret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "123",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(invalidKey)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()

	AuthMiddleware(dummyHandler)(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var resp response.ErrorResp
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	if resp.Error != "невалидный токен" {
		t.Errorf("Expected error 'невалидный токен', got '%s'", resp.Error)
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "123",
		"exp":     time.Now().Add(-time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(utils.GetJWTSecret())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()

	AuthMiddleware(dummyHandler)(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var resp response.ErrorResp
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	if resp.Error != "невалидный токен" {
		t.Errorf("Expected error 'невалидный токен', got '%s'", resp.Error)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "123",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(utils.GetJWTSecret())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()

	AuthMiddleware(dummyHandler)(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}