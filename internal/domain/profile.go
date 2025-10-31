package domain

import (
	"errors"
	"time"
)

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

type ProfileUpdate struct {
	ID        string
	FirstName string
	LastName  string
	Phone     string
	Role      string
	AvatarURL string
}

type ProfileSecurityUpdate struct {
	OldPassword string
	NewPassword string
}

var (
	ErrProfileNotFound = errors.New("profile not found")
)