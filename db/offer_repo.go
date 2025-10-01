package db

import (
	"database/sql"
	"time"

	"github.com/go-park-mail-ru/2025_2_Avrora/models"
)

type OfferRepo struct {
	db *sql.DB
}

func (r *Repo) Offer() *OfferRepo {
	return &OfferRepo{db: r.GetDB()}
}

// FindByID возвращает предложение по ID
func (or *OfferRepo) FindByID(id int) (*models.Offer, error) {
	offer := &models.Offer{}
	err := or.db.QueryRow(`
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
		return nil, err
	}
	return offer, nil
}

// FindAll возвращает все предложения с пагинацией
func (or *OfferRepo) FindAll(page, limit int) ([]models.Offer, error) {
	offset := (page - 1) * limit

	rows, err := or.db.Query(`
		SELECT id, user_id, location_id, category_id, title, description, image, price, area, rooms, address, offer_type, created_at, updated_at
		FROM offer
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offers []models.Offer
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
		offers = append(offers, o)
	}

	return offers, rows.Err()
}

// Create добавляет новое предложение
func (or *OfferRepo) Create(offer *models.Offer) error {
	offer.CreatedAt = time.Now()
	offer.UpdatedAt = offer.CreatedAt

	return or.db.QueryRow(`
		INSERT INTO offer (user_id, location_id, category_id, title, description, image, price, area, rooms, address, offer_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id
	`,
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

// Update обновляет предложение
func (or *OfferRepo) Update(offer *models.Offer) error {
	offer.UpdatedAt = time.Now()

	_, err := or.db.Exec(`
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

// CountAll возвращает количество всех предложений
func (or *OfferRepo) CountAll() (int, error) {
	var total int
	err := or.db.QueryRow("SELECT COUNT(*) FROM offer").Scan(&total)
	return total, err
}

// ClearOfferTable очищает таблицу offer
func (or *OfferRepo) ClearOfferTable() {
	_, _ = or.db.Exec("DELETE FROM offer")
}
