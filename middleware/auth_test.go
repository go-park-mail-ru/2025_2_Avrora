package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_Avrora/response"
	"github.com/go-park-mail-ru/2025_2_Avrora/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "userID not found in context", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(userID))
}

func TestAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	jwtGen := utils.NewJwtGenerator("test_secret_32_chars_min_for_tests")
	handler := AuthMiddleware(jwtGen)(http.HandlerFunc(dummyHandler))
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var resp response.ErrorResp
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	if resp.Error != "требуется авторизация" {
		t.Errorf("Expected error 'требуется авторизация', got '%s'", resp.Error)
	}
}

func TestAuthMiddleware_WrongSchema(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Token 12345")
	w := httptest.NewRecorder()

	jwtGen := utils.NewJwtGenerator("test_secret_32_chars_min_for_tests")
	handler := AuthMiddleware(jwtGen)(http.HandlerFunc(dummyHandler))
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var resp response.ErrorResp
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	if resp.Error != "некорректный формат токена" {
		t.Errorf("Expected error 'некорректный формат токена', got '%s'", resp.Error)
	}
}

func TestAuthMiddleware_EmptyToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer ") // Токен отсутствует
	w := httptest.NewRecorder()

	jwtGen := utils.NewJwtGenerator("test_secret_32_chars_min_for_tests")
	handler := AuthMiddleware(jwtGen)(http.HandlerFunc(dummyHandler))
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var resp response.ErrorResp
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	if resp.Error != "некорректный формат токена" {
    	t.Errorf("Expected error 'некорректный формат токена', got '%s'", resp.Error)
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

	jwtGen := utils.NewJwtGenerator("correct_secret")
	handler := AuthMiddleware(jwtGen)(http.HandlerFunc(dummyHandler))
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var resp response.ErrorResp
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	if resp.Error != "недействительный токен" {
		t.Errorf("Expected error 'недействительный токен', got '%s'", resp.Error)
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	jwtGen := utils.NewJwtGenerator("test_secret_32_chars_min_for_tests")
	signedString, _ := jwtGen.GenerateExpiredJWT()

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+signedString)
	w := httptest.NewRecorder()

	handler := AuthMiddleware(jwtGen)(http.HandlerFunc(dummyHandler))
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var resp response.ErrorResp
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	if resp.Error != "недействительный токен" {
		t.Errorf("Expected error 'недействительный токен', got '%s'", resp.Error)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	jwtGen := utils.NewJwtGenerator("test_secret_32_chars_min_for_tests")
	signedString, _ := jwtGen.GenerateJWT("user_123") // ← Генерируем ВАЛИДНЫЙ токен!

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+signedString)
	w := httptest.NewRecorder()

	handler := AuthMiddleware(jwtGen)(http.HandlerFunc(dummyHandler))
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Проверяем, что в ответе именно тот user_id, что в токене
	if w.Body.String() != "user_123" {
		t.Errorf("Expected user_id 'user_123', got '%s'", w.Body.String())
	}
}

func TestAuthMiddleware(t *testing.T) {
	jwtGen := utils.NewJwtGenerator("test-secret")
	token, _ := jwtGen.GenerateJWT("123")

	req := httptest.NewRequest("GET", "/api/v1/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	handler := AuthMiddleware(jwtGen)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := GetUserFromContext(r.Context())
		assert.Equal(t, "123", userID)
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}