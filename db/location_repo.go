package db

import (
	"database/sql"
	"errors"

	"github.com/go-park-mail-ru/2025_2_Avrora/models"
)

type LocationRepo struct {
	db *sql.DB
}

func (r *Repo) Location() *LocationRepo {
	return &LocationRepo{db: r.GetDB()}
}

func (lr *LocationRepo) FindByAddress(regionID int, street, houseNumber string) (*models.Location, error) {
	location := &models.Location{}
	err := lr.db.QueryRow(`
		SELECT id, region_id, street, house_number, latitude, longitude, created_at
		FROM location
		WHERE region_id = $1 AND street = $2 AND house_number = $3
	`, regionID, street, houseNumber).Scan(
		&location.ID,
		&location.RegionID,
		&location.Street,
		&location.HouseNumber,
		&location.Latitude,
		&location.Longitude,
		&location.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return location, nil
}

func (lr *LocationRepo) Create(location *models.Location) error {
	return lr.db.QueryRow(`
		INSERT INTO location (region_id, street, house_number, latitude, longitude, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id
	`,
		location.RegionID,
		location.Street,
		location.HouseNumber,
		location.Latitude,
		location.Longitude,
	).Scan(&location.ID)
}