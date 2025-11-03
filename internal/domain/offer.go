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
	ImageURLs        []string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type OfferFilter struct {
  OfferType    *string // 'sale' | 'rent'
  PropertyType *string // 'house' | 'apartment'
  Rooms        *int
  PriceMin     *int64
  PriceMax     *int64
  AreaMin      *float64
  AreaMax      *float64
  Status       *string
  Utug         *bool
  // 'active' | 'sold' | 'archived'
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
