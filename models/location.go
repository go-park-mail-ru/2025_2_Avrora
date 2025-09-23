package models

import "time"

type Location struct {
	ID          int       `json:"id"`
	RegionID    int       `json:"region_id"`
	Street      string    `json:"street"`
	HouseNumber string    `json:"house_number"`
	Latitude    float64   `json:"latitude,omitempty"`
	Longitude   float64   `json:"longitude,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}