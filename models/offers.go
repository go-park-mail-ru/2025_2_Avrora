package models

import "time"

type Offer struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Price       int       `json:"price"`
	Area        float64   `json:"area,omitempty"`
	Rooms       int       `json:"rooms,omitempty"`
	Address     string    `json:"address"`
	OfferType   string    `json:"offer_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
