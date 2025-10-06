package domain

import (
	"errors"
	"time"
)

type Offer struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	LocationID  int       `json:"location_id"`
	CategoryID  int       `json:"category_id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Image       string    `json:"image,omitempty"`
	Price       int       `json:"price"`
	Area        float64   `json:"area,omitempty"`
	Rooms       int       `json:"rooms,omitempty"`
	Address     string    `json:"address"`
	OfferType   string    `json:"offer_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrOfferNotFound = errors.New("offer not found")
)