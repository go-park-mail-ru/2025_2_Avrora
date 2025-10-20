package domain

import "time"

type User struct {
	ID        int
	Email     string
	Password  string
	CreatedAt time.Time
}

var (
	ErrUserNotFound       = Err("user not found")
	ErrInvalidCredentials = Err("invalid email or password")
)
