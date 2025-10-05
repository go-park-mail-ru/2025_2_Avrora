package utils

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher struct {
	pepper string
}

func NewPasswordHasher(pepper string) (*PasswordHasher, error) {
	if pepper == "" {
		return nil, errors.New("pepper не может быть пустым")
	}
	return &PasswordHasher{pepper: pepper}, nil
}

func (ph *PasswordHasher) Hash(password string) (string, error) {
	if password == "" {
		return "", errors.New("пароль не может быть пустым")
	}

	pepperedPassword := password + ph.pepper

	hash, err := bcrypt.GenerateFromPassword([]byte(pepperedPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("ошибка хеширования пароля: %w", err)
	}

	return string(hash), nil
}

func (ph *PasswordHasher) Compare(password, hash string) bool {
	pepperedPassword := password + ph.pepper
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pepperedPassword))
	return err == nil
}