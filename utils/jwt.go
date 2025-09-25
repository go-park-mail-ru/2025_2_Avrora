package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtGenerator struct {
	secret []byte
}

func NewJwtGenerator(secret string) *JwtGenerator {
	if secret == "" {
		panic("JWT secret cannot be empty")
	}
	return &JwtGenerator{
		secret: []byte(secret),
	}
}

func (j *JwtGenerator) GenerateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), 
	})

	return token.SignedString(j.secret)
}

func (j *JwtGenerator) GenerateExpiredJWT() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(-1 * time.Minute).Unix(),
	})
	return token.SignedString(j.secret)
}

func (j *JwtGenerator) ValidateJWT(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неподдерживаемый метод подписи")
		}
		return j.secret, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["user_id"].(string); ok {
			return userID, nil
		}
		return "", errors.New("user_id отсутствует в токене")
	}

	return "", errors.New("недействительный токен")
}