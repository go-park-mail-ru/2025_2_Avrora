package domain

//go:generate easyjson -all $GOFILE

import (
	"errors"
	"time"
)
//easyjson:json
type HousingComplex struct {
	ID            string // UUID
	Name          string
	Description   string
	YearBuilt     *int   // nullable
	LocationID    string // UUID
	Developer     string
	Address       string
	StartingPrice *int64 // nullable
	ImageURLs     []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
//easyjson:json
type ComplexInFeed struct {
	ID            string
	Name          string
	StartingPrice *int64
	Address       string
	Metro         string
	ImageURL      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
//easyjson:json
type ComplexesInFeed struct {
	Meta struct {
		Total  int
		Offset int
	}
	Complexes []ComplexInFeed
}

var (
	ErrComplexNotFound = errors.New("housing complex not found")
)
