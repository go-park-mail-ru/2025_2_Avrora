package domain
//go:generate easyjson -all $GOFILE
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
//easyjson:json
type User struct {
	ID           string
	Email        string
	PasswordHash string
	Role         UserRole
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
//easyjson:json
type UserEmailUpdate struct {
	Email string
}

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
)