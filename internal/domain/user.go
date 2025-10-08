package domain

import "time"

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	ErrUserNotFound = Err("user not found")
	ErrInvalidCredentials = Err("invalid email or password")
)