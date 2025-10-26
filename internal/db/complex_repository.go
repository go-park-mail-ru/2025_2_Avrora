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

// --- QUERIES ---

const (
	// GetByID: full complex + later load photos
	getComplexByIDQuery = `
		SELECT 
			id, name, description, year_built, location_id, developer,
			address, starting_price, created_at, updated_at
		FROM housing_complex
		WHERE id = $1`

	// List: optimized for feed (with metro, 1 image, total count)
	listComplexesInFeedQuery = `
		WITH total AS (SELECT COUNT(*) AS total_count FROM housing_complex)
		SELECT
			hc.id,
			hc.name,
			hc.starting_price,
			hc.address,
			COALESCE(
				(SELECT ms.name
				 FROM location_metro lm
				 JOIN metro_station ms ON ms.id = lm.metro_station_id
				 WHERE lm.location_id = hc.location_id
				 ORDER BY lm.distance_meters ASC
				 LIMIT 1),
				''
			) AS metro,
			COALESCE(
				(SELECT cp.url
				 FROM complex_photo cp
				 WHERE cp.complex_id = hc.id
				 ORDER BY cp.created_at ASC
				 LIMIT 1),
				''
			) AS image_url,
			hc.created_at,
			hc.updated_at,
			total.total_count
		FROM housing_complex hc
		CROSS JOIN total
		ORDER BY hc.created_at DESC
		LIMIT $1 OFFSET $2`

	createComplexQuery = `
		INSERT INTO housing_complex (
			name, description, year_built, location_id,
			developer, address, starting_price
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7
		)`

	createComplexPhotosQuery = `
		INSERT INTO complex_photo (complex_id, url, created_at, updated_at)
		SELECT $1, url, $2, $2
		FROM UNNEST($3::TEXT[]) AS url`

	updateComplexQuery = `
		UPDATE housing_complex SET
			name = $2, description = $3, year_built = $4, location_id = $5,
			developer = $6, address = $7, starting_price = $8, updated_at = $9
		WHERE id = $1`

	deleteComplexPhotosQuery = "DELETE FROM complex_photo WHERE complex_id = $1"
	deleteComplexQuery       = "DELETE FROM housing_complex WHERE id = $1"
)

type HousingComplexRepository struct {
	db  *sql.DB
	log *log.Logger
}

func NewHousingComplexRepository(db *sql.DB, log *log.Logger) *HousingComplexRepository {
	return &HousingComplexRepository{db: db, log: log}
}

// scanComplex scans a row into domain.HousingComplex (without photos)
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

// GetByID returns full complex with all photos
func (r *HousingComplexRepository) GetByID(ctx context.Context, id string) (*domain.HousingComplex, error) {
	complex, err := scanComplex(r.db.QueryRowContext(ctx, getComplexByIDQuery, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrComplexNotFound
		}
		r.log.Error(ctx, "failed to get housing complex", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	// Load all photos
	rows, err := r.db.QueryContext(ctx,
		"SELECT url FROM complex_photo WHERE complex_id = $1 ORDER BY created_at",
		id)
	if err != nil {
		r.log.Warn(ctx, "failed to load photos", zap.String("id", id), zap.Error(err))
		return complex, nil
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			continue
		}
		urls = append(urls, url)
	}
	complex.ImageURLs = urls

	return complex, nil
}

// List returns complexes in feed format with pagination metadata
func (r *HousingComplexRepository) List(ctx context.Context, page, limit int) (*domain.ComplexesInFeed, error) {
	offset := (page - 1) * limit

	rows, err := r.db.QueryContext(ctx, listComplexesInFeedQuery, limit, offset)
	if err != nil {
		r.log.Error(ctx, "failed to list complexes in feed", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var complexes []domain.ComplexInFeed
	var totalCount int

	for rows.Next() {
		var c domain.ComplexInFeed
		var total int
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.StartingPrice,
			&c.Address,
			&c.Metro,
			&c.ImageURL,
			&c.CreatedAt,
			&c.UpdatedAt,
			&total,
		)
		if err != nil {
			r.log.Error(ctx, "failed to scan complex in feed", zap.Error(err))
			return nil, err
		}
		if totalCount == 0 {
			totalCount = total
		}
		complexes = append(complexes, c)
	}

	if err = rows.Err(); err != nil {
		r.log.Error(ctx, "row iteration error", zap.Error(err))
		return nil, err
	}

	result := &domain.ComplexesInFeed{
		Complexes: complexes,
	}
	result.Meta.Total = totalCount
	result.Meta.Offset = offset

	return result, nil
}

// Create inserts a new housing complex and its photos (in app-layer transaction)
func (r *HousingComplexRepository) Create(ctx context.Context, c *domain.HousingComplex) error {
	now := time.Now().UTC()
	c.CreatedAt = now
	c.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, createComplexQuery,
		c.Name,
		c.Description,
		c.YearBuilt,
		c.LocationID,
		c.Developer,
		c.Address,
		c.StartingPrice,
	)
	if err != nil {
		r.log.Error(ctx, "failed to create housing complex", zap.Error(err))
		return err
	}

	// Insert photos in one query
	if len(c.ImageURLs) > 0 {
		_, err = r.db.ExecContext(ctx, createComplexPhotosQuery, c.ID, now, pq.StringArray(c.ImageURLs))
		if err != nil {
			r.log.Warn(ctx, "failed to insert photos", zap.String("complex_id", c.ID), zap.Error(err))
			// Note: You may want to roll back complex creation — handle in service layer with tx
		}
	}

	r.log.Info(ctx, "created housing complex", zap.String("id", c.ID))
	return nil
}

// Update modifies complex and replaces all photos
func (r *HousingComplexRepository) Update(ctx context.Context, c *domain.HousingComplex) error {
	c.UpdatedAt = time.Now().UTC()

	_, err := r.db.ExecContext(ctx, updateComplexQuery,
		c.ID,
		c.Name,
		c.Description,
		c.YearBuilt,
		c.LocationID,
		c.Developer,
		c.Address,
		c.StartingPrice,
		c.UpdatedAt,
	)
	if err != nil {
		r.log.Error(ctx, "failed to update housing complex", zap.String("id", c.ID), zap.Error(err))
		return err
	}

	// Replace photos
	_, _ = r.db.ExecContext(ctx, deleteComplexPhotosQuery, c.ID)
	if len(c.ImageURLs) > 0 {
		now := time.Now().UTC()
		_, _ = r.db.ExecContext(ctx, createComplexPhotosQuery, c.ID, now, pq.StringArray(c.ImageURLs))
	}

	r.log.Info(ctx, "updated housing complex", zap.String("id", c.ID))
	return nil
}

// Delete removes complex (photos auto-deleted via CASCADE)
func (r *HousingComplexRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, deleteComplexQuery, id)
	if err != nil {
		r.log.Error(ctx, "failed to delete housing complex", zap.String("id", id), zap.Error(err))
		return err
	}
	r.log.Info(ctx, "deleted housing complex", zap.String("id", id))
	return nil
}