package domain

import (
	"errors"
	"time"
)

type Offer struct {
	ID               int
	InHousingComplex bool
	HousingComplex   string
	OfferType        string
	PropertyType     string
	Category         string
	Address          string
	Floor            int
	TotalFloors      int
	Rooms            int
	Area             float64
	LivingArea       float64
	KitchenArea      float64
	Price            float64
	Description      string
	Deposit          float64
	Commission       float64
	RentalPeriod     string
	ImageURLs        []string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Офферы
type OfferInFeed struct {
	ID           int
	UserID       int
	OfferURL     string
	OfferType    string // sale | rent
	PropertyType string // house | apartment
	Price        float64
	Area         float64
	Rooms        int
	Floor        int
	TotalFloors  int
	Address      string
	Metro        string
	ImageURL     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Оффер в нлентеx
type OffersInFeed struct {
	Meta struct {
		Total  int
		Offset int
	}
	Offers []OfferInFeed
}
type OfferCreate struct {
	InHousingComplex bool
	HousingComplex   string
	OfferType        string
	PropertyType     string
	Category         string
	Address          string
	Floor            int
	TotalFloors      int
	Rooms            int
	Area             float64
	LivingArea       float64
	KitchenArea      float64
	Price            float64
	Description      string
	Deposit          float64
	Commission       float64
	RentalPeriod     string
	ImageURLs        []string
}

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrOfferNotFound = errors.New("offer not found")
)
