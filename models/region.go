package models

import "time"

type Region struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	ParentID  *int       `json:"parent_id,omitempty"`
	Level     int       `json:"level"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}