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
	err := r.db.QueryRow(`
		SELECT id, user_id, location_id, category_id, title, description, image, price, area, rooms, address, offer_type, created_at, updated_at
		FROM offer
		WHERE id = $1
	`, id).Scan(
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
			return &models.Offer{}, models.ErrOfferNotFound
		}
		return &models.Offer{}, err
	}
	return &offer, nil
}

func (r *OfferRepository) List(page, limit int) ([]*models.Offer, error) {
	offset := (page - 1) * limit

	rows, err := r.db.Query(`
		SELECT id, user_id, location_id, category_id, title, description, image, price, area, rooms, address, offer_type, created_at, updated_at
		FROM offer
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
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

	return r.db.QueryRow(`
		INSERT INTO offer (id, user_id, location_id, category_id, title, description, image, price, area, rooms, address, offer_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`,
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
}

func (r *OfferRepository) Update(offer *models.Offer) error {
	offer.UpdatedAt = time.Now()

	_, err := r.db.Exec(`
		UPDATE offer
		SET title = $1, description = $2, image = $3, price = $4, area = $5, rooms = $6, address = $7, offer_type = $8, updated_at = $9
		WHERE id = $10
	`,
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
	return err
}

func (r *OfferRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM offer WHERE id = $1", id)
	return err
}

func (r *OfferRepository) CountAll() (int, error) {
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM offer").Scan(&total)
	return total, err
}

func (r *OfferRepository) ListByUserID(userID string) ([]*models.Offer, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, location_id, category_id, title, description, image, price, area, rooms, address, offer_type, created_at, updated_at
		FROM offer
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
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