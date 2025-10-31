package utils

import "github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"

// Address в Location

func AddressToLocation(address string) *domain.Location {
	return &domain.Location{
		ID: "44444444-4444-4444-4444-444444444444",
		RegionID: "33333333-3333-3333-3333-333333333333",
		Latitude: 55.75580000,
		Longitude: 37.61760000,
	}
}