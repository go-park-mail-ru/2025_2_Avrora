package domain

import "time"

type Location struct {
	ID        string
	RegionID  string
	Latitude  float64
	Longitude float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GetLocation struct {
	ID        string
	Region    string // Region name for instance: 'Moscow' | 'Russia' | 'Tverskoy District' etc.
	Latitude  float64
	Longitude float64
	CreatedAt time.Time
	UpdatedAt time.Time
}