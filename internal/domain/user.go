package domain

import "time"

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRepository interface {
	Create(user User) error
	GetByID(id string) (User, error)
	GetByEmail(email string) (User, error)
}

var (
	ErrUserNotFound = Err("user not found")
	ErrInvalidCredentials = Err("invalid email or password")
)