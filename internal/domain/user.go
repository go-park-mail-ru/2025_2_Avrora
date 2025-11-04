package domain

import (
	"errors"
	"time"
)

type UserRole string

const (
	UserRoleUser    UserRole = "user"
	UserRoleOwner   UserRole = "owner"
	UserRoleRealtor UserRole = "realtor"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Role         UserRole
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserEmailUpdate struct {
	Email string
}

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
)