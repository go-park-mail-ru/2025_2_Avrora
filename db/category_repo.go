package db

import (
	"database/sql"
	"errors"
	"time"

	"github.com/go-park-mail-ru/2025_2_Avrora/models"
)

type CategoryRepo struct {
	db *sql.DB
}

func (r *Repo) Category() *CategoryRepo {
	return &CategoryRepo{db: r.GetDB()}
}

func (cr *CategoryRepo) FindBySlug(slug string) (*models.Category, error) {
	category := &models.Category{}
	err := cr.db.QueryRow(`
		SELECT id, name, slug, description, created_at
		FROM category
		WHERE slug = $1
	`, slug).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&category.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return category, nil
}

func (cr *CategoryRepo) GetAll() ([]models.Category, error) {
	rows, err := cr.db.Query(`
		SELECT id, name, slug, description, created_at
		FROM category
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Slug,
			&c.Description,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, rows.Err()
}

func (cr *CategoryRepo) Create(category *models.Category) error {
	category.CreatedAt = time.Now()

	return cr.db.QueryRow(`
		INSERT INTO category (name, slug, description, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`,
		category.Name,
		category.Slug,
		category.Description,
		category.CreatedAt,
	).Scan(&category.ID)
}