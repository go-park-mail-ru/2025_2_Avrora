package domain

import (
	"errors"
	"time"
)

type OfferType string
type PropertyType string
type OfferStatus string
type OfferID string
type PhotoURL string

const (
	OfferTypeSale OfferType = "sale"
	OfferTypeRent OfferType = "rent"

	PropertyTypeHouse     PropertyType = "house"
	PropertyTypeApartment PropertyType = "apartment"

	OfferStatusActive   OfferStatus = "active"
	OfferStatusSold     OfferStatus = "sold"
	OfferStatusArchived OfferStatus = "archived"
)

type Offer struct {
	ID               string  // UUID
	UserID           string  // UUID
	LocationID       string  // UUID
	HousingComplexID *string // UUID (nullable)
	Title            string
	Description      string
	Price            int64   // BIGINT
	Area             float64 // DECIMAL(10,2)
	Address          string
	Rooms            int
	PropertyType     PropertyType
	OfferType        OfferType
	Status           OfferStatus
	Floor            *int     // nullable
	TotalFloors      *int     // nullable
	Deposit          *int64   // nullable BIGINT
	Commission       *int64   // nullable BIGINT
	RentalPeriod     *string  // nullable
	LivingArea       *float64 // nullable
	KitchenArea      *float64 // nullable
	Metro            *string
	ImageURLs        []string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type OfferFilter struct {
	OfferType    *string  `json:"offer_type"`
	PropertyType *string  `json:"property_type"`
	Rooms        *int     `json:"rooms"`
	PriceMin     *int64   `json:"price_min"`
	PriceMax     *int64   `json:"price_max"`
	AreaMin      *float64 `json:"area_min"`
	AreaMax      *float64 `json:"area_max"`
	Status       *string  `json:"status"`
	Utug         *bool    `json:"utug"`
	Address      *string  `json:"address"`
}

// For feed (simplified + joined data)
type OfferInFeed struct {
	ID           string
	UserID       string
	OfferType    OfferType
	PropertyType PropertyType
	Price        int64
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

type PricePoint struct {
	Date  time.Time `json:"date"`
	Price int64     `json:"price"`
}

type OffersInFeed struct {
	Meta struct {
		Total  int
		Offset int
	}
	Offers []OfferInFeed
}

type OfferCreate struct {
	HousingComplexID *string
	OfferType        OfferType
	PropertyType     PropertyType
	Title            string
	Description      string
	Price            int64
	Area             float64
	Address          string
	Rooms            int
	Floor            *int
	TotalFloors      *int
	Deposit          *int64
	Commission       *int64
	RentalPeriod     *string
	LivingArea       *float64
	KitchenArea      *float64
	ImageURLs        []string
}

type FirstPhotosForOffers struct { // For offers in feed
	Photos map[OfferID]PhotoURL
}

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrOfferNotFound = errors.New("offer not found")
)
