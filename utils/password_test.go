package utils

import (
	"errors"
	"strings"
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
	ph, err := NewPasswordHasher("pepper123")
	if err != nil {
		t.Fatalf("не ожидалось ошибки при создании PasswordHasher: %v", err)
	}

	hash, err := ph.Hash("my_password")
	if err != nil {
		t.Fatalf("ошибка при хешировании: %v", err)
	}

	if hash == "" {
		t.Fatal("ожидался не пустой хеш")
	}
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	ph, err := NewPasswordHasher("pepper123")
	if err != nil {
		t.Fatalf("не ожидалось ошибки при создании PasswordHasher: %v", err)
	}

	_, err = ph.Hash("")
	EqualErrors(t, errors.New("пароль не может быть пустым"), err)
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

// Что-то типо утилиты для работы TestHashPassword_EmptyPassword
func EqualErrors(t *testing.T, expected, received error) {
	if expected == nil {
		if received != nil {
			t.Fatalf("ожидалась ошибка nil, получили %v", received)
		}
		return
	}

	if received == nil {
		t.Fatalf("ожидалась ошибка %v, получили nil", expected)
	}

	if !strings.Contains(received.Error(), expected.Error()) {
		t.Fatalf("ожидалась ошибка, содержащая %q, получили %q", expected.Error(), received.Error())
	}
}
