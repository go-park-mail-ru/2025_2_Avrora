package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	selectComplexBase = `
		id, name, description, year_built, location_id, developer,
		address, starting_price, created_at, updated_at`

	getComplexByIDQuery = `
		SELECT ` + selectComplexBase + `
		FROM housing_complex
		WHERE id = $1`

	listComplexesQuery = `
		SELECT ` + selectComplexBase + `
		FROM housing_complex
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	createComplexQuery = `
		INSERT INTO housing_complex (
			name, description, year_built, location_id, developer,
			address, starting_price, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	updateComplexQuery = `
		UPDATE housing_complex SET
			name = $1, description = $2, year_built = $3, location_id = $4,
			developer = $5, address = $6, starting_price = $7, updated_at = $8
		WHERE id = $9`

	deleteComplexQuery     = "DELETE FROM housing_complex WHERE id = $1"
	countAllComplexesQuery = "SELECT COUNT(*) FROM housing_complex"
)

type HousingComplexRepository struct {
	db  *sql.DB
	log *log.Logger
}

func NewHousingComplexRepository(db *sql.DB, log *log.Logger) *HousingComplexRepository {
	return &HousingComplexRepository{db: db, log: log}
}

// scanComplex scans a single row into a domain.HousingComplex (without photos)
func scanComplex(row *sql.Row) (*domain.HousingComplex, error) {
	var c domain.HousingComplex
	var yearBuilt *int
	var startingPrice *int64

	err := row.Scan(
		&c.ID,
		&c.Name,
		&c.Description,
		&yearBuilt,
		&c.LocationID,
		&c.Developer,
		&c.Address,
		&startingPrice,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	c.YearBuilt = yearBuilt
	c.StartingPrice = startingPrice
	return &c, nil
}

// scanComplexes scans multiple rows
func scanComplexes(rows *sql.Rows) ([]*domain.HousingComplex, error) {
	var complexes []*domain.HousingComplex
	for rows.Next() {
		var c domain.HousingComplex
		var yearBuilt *int
		var startingPrice *int64

		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Description,
			&yearBuilt,
			&c.LocationID,
			&c.Developer,
			&c.Address,
			&startingPrice,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		c.YearBuilt = yearBuilt
		c.StartingPrice = startingPrice
		complexes = append(complexes, &c)
	}
	return complexes, nil
}

// fetchPhotosForComplexes loads photos from complex_photo table
func (r *HousingComplexRepository) fetchPhotosForComplexes(ctx context.Context, complexes []*domain.HousingComplex) error {
	if len(complexes) == 0 {
		return nil
	}

	ids := make([]string, len(complexes))
	for i, c := range complexes {
		ids[i] = c.ID
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT complex_id, url FROM complex_photo WHERE complex_id = ANY($1) ORDER BY created_at`,
		pq.Array(ids))
	if err != nil {
		return err
	}
	defer rows.Close()

	photoMap := make(map[string][]string)
	for rows.Next() {
		var complexID, url string
		if err := rows.Scan(&complexID, &url); err != nil {
			return err
		}
		photoMap[complexID] = append(photoMap[complexID], url)
	}

	for _, c := range complexes {
		c.ImageURLs = photoMap[c.ID]
	}

	return nil
}

// GetByID fetches a housing complex by ID with its photos
func (r *HousingComplexRepository) GetByID(ctx context.Context, id string) (*domain.HousingComplex, error) {
	complex, err := scanComplex(r.db.QueryRowContext(ctx, getComplexByIDQuery, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrComplexNotFound
		}
		r.log.Error(ctx, "failed to get housing complex", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	if err := r.fetchPhotosForComplexes(ctx, []*domain.HousingComplex{complex}); err != nil {
		r.log.Warn(ctx, "failed to load photos for complex", zap.String("id", id), zap.Error(err))
	}

	return complex, nil
}

// List returns paginated housing complexes with photos
func (r *HousingComplexRepository) List(ctx context.Context, page, limit int) ([]*domain.HousingComplex, error) {
	offset := (page - 1) * limit
	rows, err := r.db.QueryContext(ctx, listComplexesQuery, limit, offset)
	if err != nil {
		r.log.Error(ctx, "failed to list housing complexes", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	complexes, err := scanComplexes(rows)
	if err != nil {
		r.log.Error(ctx, "failed to scan housing complexes", zap.Error(err))
		return nil, err
	}

	if err := r.fetchPhotosForComplexes(ctx, complexes); err != nil {
		r.log.Warn(ctx, "partial photo load for complexes", zap.Error(err))
	}

	return complexes, nil
}

// Create inserts a new housing complex and its photos
func (r *HousingComplexRepository) Create(ctx context.Context, c *domain.HousingComplex) error {
	now := time.Now().UTC()
	c.CreatedAt = now
	c.UpdatedAt = now

	err := r.db.QueryRowContext(ctx, createComplexQuery,
		c.Name,
		c.Description,
		c.YearBuilt,
		c.LocationID,
		c.Developer,
		c.Address,
		c.StartingPrice,
		now,
		now,
	).Scan(&c.ID)
	if err != nil {
		r.log.Error(ctx, "failed to create housing complex", zap.Error(err))
		return err
	}

	// Insert photos
	for _, url := range c.ImageURLs {
		_, err := r.db.ExecContext(ctx,
			"INSERT INTO complex_photo (complex_id, url, created_at, updated_at) VALUES ($1, $2, $3, $3)",
			c.ID, url, now)
		if err != nil {
			r.log.Warn(ctx, "failed to insert complex photo", zap.String("complex_id", c.ID), zap.String("url", url), zap.Error(err))
		}
	}

	r.log.Info(ctx, "created housing complex", zap.String("id", c.ID))
	return nil
}

// Update modifies an existing housing complex and replaces its photos
func (r *HousingComplexRepository) Update(ctx context.Context, c *domain.HousingComplex) error {
	c.UpdatedAt = time.Now().UTC()

	_, err := r.db.ExecContext(ctx, updateComplexQuery,
		c.Name,
		c.Description,
		c.YearBuilt,
		c.LocationID,
		c.Developer,
		c.Address,
		c.StartingPrice,
		c.UpdatedAt,
		c.ID,
	)
	if err != nil {
		r.log.Error(ctx, "failed to update housing complex", zap.String("id", c.ID), zap.Error(err))
		return err
	}

	// Replace photos: delete all and reinsert
	_, _ = r.db.ExecContext(ctx, "DELETE FROM complex_photo WHERE complex_id = $1", c.ID)
	now := time.Now().UTC()
	for _, url := range c.ImageURLs {
		_, _ = r.db.ExecContext(ctx,
			"INSERT INTO complex_photo (complex_id, url, created_at, updated_at) VALUES ($1, $2, $3, $3)",
			c.ID, url, now)
	}

	r.log.Info(ctx, "updated housing complex", zap.String("id", c.ID))
	return nil
}

// Delete removes a housing complex (photos deleted via CASCADE)
func (r *HousingComplexRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, deleteComplexQuery, id)
	if err != nil {
		r.log.Error(ctx, "failed to delete housing complex", zap.String("id", id), zap.Error(err))
		return err
	}
	r.log.Info(ctx, "deleted housing complex", zap.String("id", id))
	return nil
}