package handlers

import "errors"

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	Email string       `json:"email"`
}

type RegisterResponse struct {
	Email string `json:"email"`
}

type LogoutResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

var (
	ErrInvalidEmail = errors.New("неправильный формат почты")
	ErrInvalidJSON  = errors.New("неправильный формат JSON")
)