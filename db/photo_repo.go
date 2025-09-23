package db

import (
	"database/sql"
	"time"

	"github.com/go-park-mail-ru/2025_2_Avrora/models"
)

type PhotoRepo struct {
	db *sql.DB
}

func (r *Repo) Photo() *PhotoRepo {
	return &PhotoRepo{db: r.GetDB()}
}

func (pr *PhotoRepo) GetByOfferID(offerID int) ([]models.Photo, error) {
	rows, err := pr.db.Query(`
		SELECT id, offer_id, url, position, uploaded_at
		FROM photo
		WHERE offer_id = $1
		ORDER BY position
	`, offerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []models.Photo
	for rows.Next() {
		var p models.Photo
		err := rows.Scan(
			&p.ID,
			&p.OfferID,
			&p.URL,
			&p.Position,
			&p.UploadedAt,
		)
		if err != nil {
			return nil, err
		}
		photos = append(photos, p)
	}

	return photos, rows.Err()
}

func (pr *PhotoRepo) Create(photo *models.Photo) error {
	photo.UploadedAt = time.Now()

	return pr.db.QueryRow(`
		INSERT INTO photo (offer_id, url, position, uploaded_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`,
		photo.OfferID,
		photo.URL,
		photo.Position,
		photo.UploadedAt,
	).Scan(&photo.ID)
}

func (pr *PhotoRepo) DeleteByOfferID(offerID int) error {
	_, err := pr.db.Exec("DELETE FROM photo WHERE offer_id = $1", offerID)
	return err
}