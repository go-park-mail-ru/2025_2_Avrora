package domain

//go:generate easyjson -all $GOFILE
import (
	"errors"
	"time"
)

//easyjson:json
type Profile struct {
	ID        string
	UserID    string
	FirstName string
	LastName  string
	Phone     string
	Role      string
	Email     string
	AvatarURL string
	CreatedAt time.Time
	UpdatedAt time.Time
}

//easyjson:json
type ProfileUpdate struct {
	ID        string
	FirstName string
	LastName  string
	Phone     string
	Role      string
	AvatarURL string
}

//easyjson:json
type ProfileSecurityUpdate struct {
	OldPassword string
	NewPassword string
}

var (
	ErrProfileNotFound = errors.New("profile not found")
)
