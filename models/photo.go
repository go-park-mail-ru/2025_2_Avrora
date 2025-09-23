package models

import "time"

type Photo struct {
	ID         int       `json:"id"`
	OfferID    int       `json:"offer_id"`
	URL        string    `json:"url"`
	Position   int       `json:"position"`
	UploadedAt time.Time `json:"uploaded_at"`
}