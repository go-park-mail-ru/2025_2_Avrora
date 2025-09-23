package db

import (
	"database/sql"
	"errors"
	"time"

	"github.com/go-park-mail-ru/2025_2_Avrora/models"
)

type RegionRepo struct {
	db *sql.DB
}

func (r *Repo) Region() *RegionRepo {
	return &RegionRepo{db: r.GetDB()}
}

func (rr *RegionRepo) FindBySlug(slug string) (*models.Region, error) {
	region := &models.Region{}
	err := rr.db.QueryRow(`
		SELECT id, name, parent_id, level, slug, created_at
		FROM region
		WHERE slug = $1
	`, slug).Scan(
		&region.ID,
		&region.Name,
		&region.ParentID,
		&region.Level,
		&region.Slug,
		&region.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return region, nil
}

func (rr *RegionRepo) FindByID(id int) (*models.Region, error) {
	region := &models.Region{}
	err := rr.db.QueryRow(`
		SELECT id, name, parent_id, level, slug, created_at
		FROM region
		WHERE id = $1
	`, id).Scan(
		&region.ID,
		&region.Name,
		&region.ParentID,
		&region.Level,
		&region.Slug,
		&region.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return region, nil
}

func (rr *RegionRepo) GetChildren(parentID int) ([]models.Region, error) {
	rows, err := rr.db.Query(`
		SELECT id, name, parent_id, level, slug, created_at
		FROM region
		WHERE parent_id = $1
		ORDER BY name
	`, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var regions []models.Region
	for rows.Next() {
		var r models.Region
		err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.ParentID,
			&r.Level,
			&r.Slug,
			&r.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		regions = append(regions, r)
	}

	return regions, rows.Err()
}

func (rr *RegionRepo) Create(region *models.Region) error {
	region.CreatedAt = time.Now()

	return rr.db.QueryRow(`
		INSERT INTO region (name, parent_id, level, slug, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`,
		region.Name,
		region.ParentID,
		region.Level,
		region.Slug,
		region.CreatedAt,
	).Scan(&region.ID)
}