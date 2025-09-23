package utils

import (
	"os"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
    LoadEnv()

    pepper := os.Getenv("PASSWORD_PEPPER")
    if pepper == "" {
		return "", bcrypt.ErrHashTooShort
	}

    data := []byte(pepper + password)
	hash, err := bcrypt.GenerateFromPassword(data, bcrypt.DefaultCost)
	return string(hash), err
}

func CheckPasswordHash(password, hash string) bool {
	pepper := os.Getenv("PASSWORD_PEPPER")
	if pepper == "" {
		return false
	}

	data := []byte(pepper + password)
	err := bcrypt.CompareHashAndPassword([]byte(hash), data)
	return err == nil
}