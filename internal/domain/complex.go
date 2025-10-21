package domain

import "time"

type Complex struct {
	ID            int
	Description   string
	Name          string
	Address       string
	Metro         string
	Developer     string
	BuiltYear     int
	ImageURL      []string
	StartingPrice float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// комплекс в ленте
type ComplexInFeed struct {
	ID            int
	Name          string
	StartingPrice float64
	Address       string
	Metro         string
	ImageURL      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ComplexesInFeed struct {
	Complexes []ComplexInFeed
}
