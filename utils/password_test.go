package utils

import (
	"testing"
)

// -------------------
// Тест создания PasswordHasher
// -------------------
func TestNewPasswordHasher(t *testing.T) {
	ph, err := NewPasswordHasher("pepper123")
	if err != nil {
		t.Fatalf("не ожидалось ошибки: %v", err)
	}
	if ph == nil {
		t.Fatal("PasswordHasher не должен быть nil")
	}
}

func TestNewPasswordHasher_EmptyPepper(t *testing.T) {
	ph, err := NewPasswordHasher("")
	if err == nil {
		t.Fatal("ожидалась ошибка для пустого pepper")
	}
	if ph != nil {
		t.Fatal("PasswordHasher должен быть nil при пустом pepper")
	}
}

// -------------------
// Тест хеширования пароля
// -------------------
func TestHashPassword(t *testing.T) {
	ph, _ := NewPasswordHasher("pepper123")

	hash, err := ph.Hash("my_password")
	if err != nil {
		t.Fatalf("ошибка при хешировании: %v", err)
	}

	if hash == "" {
		t.Fatal("ожидался не пустой хеш")
	}
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	ph, _ := NewPasswordHasher("pepper123")

	_, err := ph.Hash("")
	if err == nil {
		t.Fatal("ожидалась ошибка при пустом пароле")
	}
}

// -------------------
// Тест сравнения пароля
// -------------------
func TestComparePassword(t *testing.T) {
	ph, _ := NewPasswordHasher("pepper123")

	password := "my_password"
	hash, _ := ph.Hash(password)

	if !ph.Compare(password, hash) {
		t.Fatal("ожидалось, что пароль совпадает с хешем")
	}

	if ph.Compare("wrong_password", hash) {
		t.Fatal("ожидалось, что неверный пароль не совпадает с хешем")
	}
}
