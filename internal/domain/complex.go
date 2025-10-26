package domain

import (
	"errors"
	"time"
)

type HousingComplex struct {
	ID            string    // UUID
	Name          string
	Description   string
	YearBuilt     *int      // nullable
	LocationID    string    // UUID
	Developer     string
	Address       string
	StartingPrice *int64	// nullable
	ImageURLs     []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

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