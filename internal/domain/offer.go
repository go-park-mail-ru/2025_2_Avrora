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

// Офферы
type OfferInFeed struct {
	ID           int     `json:"id"`
	UserID       int     `json:"user_id"`
	OfferURL     string  `json:"offer_url"`
	OfferType    string  `json:"offer_type"`    // sale | rent
	PropertyType string  `json:"property_type"` // flat | house
	Price        float64 `json:"price"`
	Area         float64 `json:"area"`
	Rooms        int     `json:"rooms"`
	Floor        int     `json:"floor"`
	TotalFloors  int     `json:"total_floors"`
	Address      string  `json:"address"`
	Metro        string  `json:"metro"`
	ImageURL     string  `json:"image_url"`
}

// Оффер в нлентеx
type OffersInFeed struct {
	Meta struct {
		Total  int `json:"total"`
		Offset int `json:"offset"`
	} `json:"meta"`
	Offers []OfferInFeed `json:"offers"`
}
type OfferCreate struct {
	InHousingComplex bool     `json:"in_housing_complex"`
	HousingComplex   string   `json:"housing_complex,omitempty"`
	OfferType        string   `json:"offer_type"`
	PropertyType     string   `json:"property_type"`
	Category         string   `json:"category"`
	Address          string   `json:"address"`
	Floor            int      `json:"floor"`
	TotalFloors      int      `json:"total_floors"`
	Rooms            int      `json:"rooms"`
	Area             float64  `json:"area"`
	LivingArea       float64  `json:"living_area"`
	KitchenArea      float64  `json:"kitchen_area"`
	Price            float64  `json:"price"`
	Description      string   `json:"description,omitempty"`
	Deposit          float64  `json:"deposit,omitempty"`
	Commission       float64  `json:"commission,omitempty"`
	RentalPeriod     string   `json:"rental_period,omitempty"`
	ImageURLs        []string `json:"image_urls,omitempty"`
}

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrOfferNotFound = errors.New("offer not found")
)
