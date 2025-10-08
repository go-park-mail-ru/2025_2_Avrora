package utils

import (
	jwt "github.com/golang-jwt/jwt/v5"
	"testing"
	"time"
)

func TestNewJwtGenerator(t *testing.T) {
	j := NewJwtGenerator("secret123")
	if j == nil {
		t.Fatal("JwtGenerator не должен быть nil")
	}
}

func TestNewJwtGenerator_EmptySecret(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("ожидался panic при пустом секрете")
		}
	}()
	NewJwtGenerator("")
}

// -------------------
// Тест генерации и проверки токена
// -------------------
func TestGenerateAndValidateJWT(t *testing.T) {
	j := NewJwtGenerator("test_secret")

	tokenStr, err := j.GenerateJWT("user_123")
	if err != nil {
		t.Fatalf("ошибка при генерации токена: %v", err)
	}

	userID, err := j.ValidateJWT(tokenStr)
	if err != nil {
		t.Fatalf("ошибка при валидации токена: %v", err)
	}

	if userID != "user_123" {
		t.Errorf("expected user_123, got %s", userID)
	}
}

// -------------------
// Тест просроченного токена
// -------------------
func TestExpiredJWT(t *testing.T) {
	j := NewJwtGenerator("test_secret")

	tokenStr, err := j.GenerateExpiredJWT()
	if err != nil {
		t.Fatalf("ошибка при генерации просроченного токена: %v", err)
	}

	_, err = j.ValidateJWT(tokenStr)
	if err == nil {
		t.Fatal("ожидалась ошибка для просроченного токена")
	}
}

// -------------------
// Тест токена с неверным секретом
// -------------------
func TestInvalidSignature(t *testing.T) {
	j1 := NewJwtGenerator("secret1")
	j2 := NewJwtGenerator("secret2")

	tokenStr, _ := j1.GenerateJWT("user_123")

	_, err := j2.ValidateJWT(tokenStr)
	if err == nil {
		t.Fatal("ожидалась ошибка для токена с неверным секретом")
	}
}

// -------------------
// Тест токена без user_id
// -------------------
func TestTokenWithoutUserID(t *testing.T) {
	j := NewJwtGenerator("test_secret")

	// Создаем токен без user_id
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenStr, _ := token.SignedString(j.secret)

	_, err := j.ValidateJWT(tokenStr)
	if err == nil {
		t.Fatal("ожидалась ошибка для токена без user_id")
	}
}
