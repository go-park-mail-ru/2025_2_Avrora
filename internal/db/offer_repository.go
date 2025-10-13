package db

import (
	"database/sql"
	"errors"
	"time"

	models "github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
)

type OfferRepository struct {
	db *sql.DB
}

func NewOfferRepository(db *sql.DB) *OfferRepository {
	return &OfferRepository{db: db}
}

func (r *OfferRepository) GetByID(id string) (*models.Offer, error) {
	offer := models.Offer{}
	err := r.db.QueryRow(getOfferByIDQuery, id).Scan(
		&offer.ID,
		&offer.UserID,
		&offer.LocationID,
		&offer.CategoryID,
		&offer.Title,
		&offer.Description,
		&offer.Image,
		&offer.Price,
		&offer.Area,
		&offer.Rooms,
		&offer.Address,
		&offer.OfferType,
		&offer.CreatedAt,
		&offer.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrOfferNotFound
		}
		return nil, err
	}
	return &offer, nil
}

func (r *OfferRepository) List(page, limit int) ([]*models.Offer, error) {
	offset := (page - 1) * limit

	rows, err := r.db.Query(listOffersQuery, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offers []*models.Offer
	for rows.Next() {
		var o models.Offer
		err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.LocationID,
			&o.CategoryID,
			&o.Title,
			&o.Description,
			&o.Image,
			&o.Price,
			&o.Area,
			&o.Rooms,
			&o.Address,
			&o.OfferType,
			&o.CreatedAt,
			&o.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		offers = append(offers, &o)
	}

	return offers, rows.Err()
}

func (r *OfferRepository) Create(offer *models.Offer) error {
	now := time.Now()
	offer.CreatedAt = now
	offer.UpdatedAt = now

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	err = tx.QueryRow(
		createOfferQuery,
		offer.ID,
		offer.UserID,
		offer.LocationID,
		offer.CategoryID,
		offer.Title,
		offer.Description,
		offer.Image,
		offer.Price,
		offer.Area,
		offer.Rooms,
		offer.Address,
		offer.OfferType,
		offer.CreatedAt,
		offer.UpdatedAt,
	).Scan(&offer.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *OfferRepository) Update(offer *models.Offer) error {
	offer.UpdatedAt = time.Now()

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		updateOfferQuery,
		offer.Title,
		offer.Description,
		offer.Image,
		offer.Price,
		offer.Area,
		offer.Rooms,
		offer.Address,
		offer.OfferType,
		offer.UpdatedAt,
		offer.ID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *OfferRepository) Delete(id string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(deleteOfferQuery, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *OfferRepository) CountAll() (int, error) {
	var total int
	err := r.db.QueryRow(countAllOffersQuery).Scan(&total)
	return total, err
}

func (r *OfferRepository) ListByUserID(userID string) ([]*models.Offer, error) {
	rows, err := r.db.Query(listOffersByUserIDQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offers []*models.Offer
	for rows.Next() {
		var o models.Offer
		err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.LocationID,
			&o.CategoryID,
			&o.Title,
			&o.Description,
			&o.Image,
			&o.Price,
			&o.Area,
			&o.Rooms,
			&o.Address,
			&o.OfferType,
			&o.CreatedAt,
			&o.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		offers = append(offers, &o)
	}

	return offers, rows.Err()
}