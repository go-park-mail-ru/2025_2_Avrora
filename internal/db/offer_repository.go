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

// SQL constants using actual schema
const (
	listPhotosForOfferQuery = `
		SELECT image_url
		FROM offer_photo
		WHERE offer_id = $1
		ORDER BY created_at DESC`

	getOfferByIDQuery = `
		SELECT
			o.id,
			o.user_id,
			o.location_id,
			o.housing_complex_id,
			o.title,
			o.description,
			o.price,
			o.area,
			o.address,
			o.rooms,
			o.property_type,
			o.offer_type,
			o.status,
			o.floor,
			o.total_floors,
			o.deposit,
			o.commission,
			o.rental_period,
			o.living_area,
			o.kitchen_area,
			COALESCE(ARRAY_AGG(op.url) FILTER (WHERE op.url IS NOT NULL), '{}') AS image_urls,
			o.created_at,
			o.updated_at
		FROM offer o
		LEFT JOIN offer_photo op ON op.offer_id = o.id
		WHERE o.id = $1
		GROUP BY
			o.id,
			o.user_id,
			o.location_id,
			o.housing_complex_id,
			o.title,
			o.description,
			o.price,
			o.area,
			o.address,
			o.rooms,
			o.property_type,
			o.offer_type,
			o.status,
			o.floor,
			o.total_floors,
			o.deposit,
			o.commission,
			o.rental_period,
			o.living_area,
			o.kitchen_area,
			o.created_at,
			o.updated_at
	`

	createOfferQuery = `
		INSERT INTO offer (
			id,
			user_id,
			location_id,
			housing_complex_id,
			title,
			description,
			price,
			area,
			address,
			rooms,
			property_type,
			offer_type,
			status,
			floor,
			total_floors,
			deposit,
			commission,
			rental_period,
			living_area,
			kitchen_area,
			created_at,
			updated_at
		) VALUES (
			$1,  -- id (UUID)
			$2,  -- user_id
			$3,  -- location_id
			$4,  -- housing_complex_id (can be NULL)
			$5,  -- title
			$6,  -- description
			$7,  -- price
			$8,  -- area
			$9,  -- address
			$10, -- rooms
			$11, -- property_type
			$12, -- offer_type
			'active', -- status (default to 'active')
			$13, -- floor (can be NULL)
			$14, -- total_floors (can be NULL)
			$15, -- deposit (can be NULL)
			$16, -- commission (can be NULL)
			$17, -- rental_period (can be NULL)
			$18, -- living_area (can be NULL)
			$19, -- kitchen_area (can be NULL)
			NOW(),
			NOW()
		)
		RETURNING id, created_at, updated_at
	`

	updateOfferQuery = `
		UPDATE offer SET
			location_id = $1, housing_complex_id = $2, title = $3, description = $4,
			price = $5, area = $6, address = $7, rooms = $8, property_type = $9,
			offer_type = $10, status = $11, floor = $12, total_floors = $13,
			deposit = $14, commission = $15, rental_period = $16,
			living_area = $17, kitchen_area = $18, updated_at = $19
		WHERE id = $20`

	deleteOfferQuery = "DELETE FROM offer WHERE id = $1"

	countAllOffersQuery = "SELECT COUNT(*) FROM offer WHERE status = 'active'"

	listOffersQuery = `
		SELECT 
			o.id,
			o.user_id,
			o.offer_type,
			o.property_type,
			o.price,
			o.area,
			o.rooms,
			o.floor,
			o.total_floors,
			o.address,
			ms.name AS metro,
			op.url AS image_url,
			o.created_at,
			o.updated_at
		FROM offer o
		LEFT JOIN (
			SELECT DISTINCT ON (location_id)
				location_id,
				metro_station_id
			FROM location_metro
			ORDER BY location_id, distance_meters ASC
		) lm ON lm.location_id = o.location_id
		LEFT JOIN metro_station ms ON ms.id = lm.metro_station_id
		LEFT JOIN (
			SELECT DISTINCT ON (offer_id)
				offer_id,
				url
			FROM offer_photo
			ORDER BY offer_id, created_at ASC
		) op ON op.offer_id = o.id
		WHERE o.status = 'active'
		ORDER BY o.created_at DESC
		LIMIT $1 OFFSET $2
	`

	listOffersByUserIDQuery = `
		SELECT
			o.id,
			o.user_id,
			o.offer_type,
			o.property_type,
			o.price,
			o.area,
			o.rooms,
			o.floor,
			o.total_floors,
			o.address,
			ms.name AS metro,
			op.url AS image_url,
			o.created_at,
			o.updated_at
		FROM offer o
		LEFT JOIN (
			SELECT DISTINCT ON (location_id)
				location_id,
				metro_station_id
			FROM location_metro
			ORDER BY location_id, distance_meters ASC
		) lm ON lm.location_id = o.location_id
		LEFT JOIN metro_station ms ON ms.id = lm.metro_station_id
		LEFT JOIN (
			SELECT DISTINCT ON (offer_id)
				offer_id,
				url
			FROM offer_photo
			ORDER BY offer_id, created_at ASC
		) op ON op.offer_id = o.id
		WHERE o.user_id = $1
		AND o.status = 'active'
		ORDER BY o.created_at DESC
		LIMIT $2 OFFSET $3
		`

	countOffersByUserIDQuery = `
		SELECT COUNT(*)
		FROM offer
		WHERE user_id = $1 AND status = 'active'
		`
)

type OfferRepository struct {
	db  *sql.DB
	log *log.Logger
}

func NewOfferRepository(db *sql.DB, log *log.Logger) *OfferRepository {
	return &OfferRepository{db: db, log: log}
}

func scanOfferRow(scanner interface {
	Scan(dest ...any) error
}) (*domain.Offer, error) {
	var (
		housingComplexID *string
		floor, totalFloors *int
		deposit, commission *int64
		rentalPeriod *string
		livingArea, kitchenArea *float64
		offer domain.Offer
	)

	err := scanner.Scan(
		&offer.ID,
		&offer.UserID,
		&offer.LocationID,
		&housingComplexID,
		&offer.Title,
		&offer.Description,
		&offer.Price,
		&offer.Area,
		&offer.Address,
		&offer.Rooms,
		&offer.PropertyType,
		&offer.OfferType,
		&offer.Status,
		&floor,
		&totalFloors,
		&deposit,
		&commission,
		&rentalPeriod,
		&livingArea,
		&kitchenArea,
		&offer.CreatedAt,
		&offer.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	offer.HousingComplexID = housingComplexID
	offer.Floor = floor
	offer.TotalFloors = totalFloors
	offer.Deposit = deposit
	offer.Commission = commission
	offer.RentalPeriod = rentalPeriod
	offer.LivingArea = livingArea
	offer.KitchenArea = kitchenArea

	return &offer, nil
}

func scanOffer(row *sql.Row) (*domain.Offer, error) {
	return scanOfferRow(row)
}

func scanOfferInFeedRow(scanner interface {
	Scan(dest ...any) error
}) (*domain.OfferInFeed, error) {
	var o domain.OfferInFeed
	var metro, imageURL *string

	err := scanner.Scan(
		&o.ID,
		&o.UserID,
		&o.OfferType,
		&o.PropertyType,
		&o.Price,
		&o.Area,
		&o.Rooms,
		&o.Floor,
		&o.TotalFloors,
		&o.Address,
		&metro,
		&imageURL,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Handle nullable fields
	if metro != nil {
		o.Metro = *metro
	}
	if imageURL != nil {
		o.ImageURL = *imageURL
	}

	return &o, nil
}

func scanOffersInFeed(rows *sql.Rows) ([]domain.OfferInFeed, error) {
	var offers []domain.OfferInFeed
	for rows.Next() {
		offer, err := scanOfferInFeedRow(rows)
		if err != nil {
			return nil, err
		}
		offers = append(offers, *offer)
	}
	return offers, rows.Err()
}

func (r *OfferRepository) fetchPhotosForOffers(ctx context.Context, offers []*domain.Offer) error {
	if len(offers) == 0 {
		return nil
	}

	ids := make([]string, len(offers))
	for i, o := range offers {
		ids[i] = o.ID
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT offer_id, url 
		FROM offer_photo 
		WHERE offer_id = ANY($1)
		ORDER BY created_at ASC
	`, pq.Array(ids))
	if err != nil {
		return err
	}
	defer rows.Close()

	// Group by offer_id
	photoMap := make(map[string][]string)
	for rows.Next() {
		var offerID, url string
		if err := rows.Scan(&offerID, &url); err != nil {
			return err
		}
		photoMap[offerID] = append(photoMap[offerID], url)
	}

	// Assign back
	for _, offer := range offers {
		offer.ImageURLs = photoMap[offer.ID]
		if offer.ImageURLs == nil {
			offer.ImageURLs = []string{}
		}
	}

	return nil
}

func (r *OfferRepository) GetByID(ctx context.Context, id string) (*domain.Offer, error) {
	offer, err := scanOffer(r.db.QueryRowContext(ctx, getOfferByIDQuery, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrOfferNotFound
		}
		r.log.Error(ctx, "failed to get offer", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	// Load photos in batch (even for 1 offer)
	if err := r.fetchPhotosForOffers(ctx, []*domain.Offer{offer}); err != nil {
		r.log.Warn(ctx, "failed to load photos for offer", zap.String("id", id), zap.Error(err))
		offer.ImageURLs = []string{}
	}

	return offer, nil
}

func (r *OfferRepository) List(ctx context.Context, page, limit int) (*domain.OffersInFeed, error) {
	offset := (page - 1) * limit

	rows, err := r.db.QueryContext(ctx, listOffersQuery, limit, offset)
	if err != nil {
		r.log.Error(ctx, "failed to list offers", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	offers, err := scanOffersInFeed(rows)
	if err != nil {
		r.log.Error(ctx, "failed to scan offers for feed", zap.Error(err))
		return nil, err
	}

	total, err := r.CountAll(ctx)
	if err != nil {
		r.log.Error(ctx, "failed to count offers", zap.Error(err))
		return nil, err
	}

	return &domain.OffersInFeed{
		Meta: struct {
			Total  int
			Offset int
		}{
			Total:  total,
			Offset: offset,
		},
		Offers: offers,
	}, nil
}

func (r *OfferRepository) Create(ctx context.Context, offer *domain.Offer) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now().UTC()
	offer.CreatedAt = now
	offer.UpdatedAt = now

	err = tx.QueryRowContext(ctx, createOfferQuery,
		offer.ID, // assume UUID generated in service layer
		offer.UserID,
		offer.LocationID,
		offer.HousingComplexID,
		offer.Title,
		offer.Description,
		offer.Price,
		offer.Area,
		offer.Address,
		offer.Rooms,
		offer.PropertyType,
		offer.OfferType,
		offer.Floor,
		offer.TotalFloors,
		offer.Deposit,
		offer.Commission,
		offer.RentalPeriod,
		offer.LivingArea,
		offer.KitchenArea,
		now,
		now,
	).Scan(&offer.ID)
	if err != nil {
		r.log.Error(ctx, "failed to create offer", zap.Error(err))
		return err
	}

	// Insert photos
	for _, url := range offer.ImageURLs {
		_, err := tx.ExecContext(ctx,
			"INSERT INTO offer_photo (offer_id, url, created_at, updated_at) VALUES ($1, $2, $3, $3)",
			offer.ID, url, now)
		if err != nil {
			r.log.Warn(ctx, "failed to insert photo", zap.String("offer_id", offer.ID), zap.String("url", url), zap.Error(err))
			// Optionally: abort on photo error? Usually not.
		}
	}

	if err := tx.Commit(); err != nil {
		r.log.Error(ctx, "failed to commit offer creation", zap.Error(err))
		return err
	}

	r.log.Info(ctx, "created offer", zap.String("id", offer.ID))
	return nil
}

func (r *OfferRepository) Update(ctx context.Context, offer *domain.Offer) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	offer.UpdatedAt = time.Now().UTC()

	_, err = tx.ExecContext(ctx, updateOfferQuery,
		offer.LocationID,
		offer.HousingComplexID,
		offer.Title,
		offer.Description,
		offer.Price,
		offer.Area,
		offer.Address,
		offer.Rooms,
		offer.PropertyType,
		offer.OfferType,
		offer.Status,
		offer.Floor,
		offer.TotalFloors,
		offer.Deposit,
		offer.Commission,
		offer.RentalPeriod,
		offer.LivingArea,
		offer.KitchenArea,
		offer.UpdatedAt,
		offer.ID,
	)
	if err != nil {
		r.log.Error(ctx, "failed to update offer", zap.String("id", offer.ID), zap.Error(err))
		return err
	}

	// Replace photos
	_, _ = tx.ExecContext(ctx, "DELETE FROM offer_photo WHERE offer_id = $1", offer.ID)
	now := time.Now().UTC()
	for _, url := range offer.ImageURLs {
		_, _ = tx.ExecContext(ctx,
			"INSERT INTO offer_photo (offer_id, url, created_at, updated_at) VALUES ($1, $2, $3, $3)",
			offer.ID, url, now)
	}

	if err := tx.Commit(); err != nil {
		r.log.Error(ctx, "failed to commit offer update", zap.Error(err))
		return err
	}

	r.log.Info(ctx, "updated offer", zap.String("id", offer.ID))
	return nil
}

func (r *OfferRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, deleteOfferQuery, id)
	if err != nil {
		r.log.Error(ctx, "failed to delete offer", zap.String("id", id), zap.Error(err))
		return err
	}
	r.log.Info(ctx, "deleted offer", zap.String("id", id))
	return nil
}

func (r *OfferRepository) CountAll(ctx context.Context) (int, error) {
	var total int
	err := r.db.QueryRowContext(ctx, countAllOffersQuery).Scan(&total)
	if err != nil {
		r.log.Error(ctx, "failed to count offers", zap.Error(err))
		return 0, err
	}
	return total, nil
}

func (r *OfferRepository) ListByUserID(ctx context.Context, userID string, page, limit int) (*domain.OffersInFeed, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10 // default or enforce min
	}
	offset := (page - 1) * limit

	// Fetch offers
	rows, err := r.db.QueryContext(ctx, listOffersByUserIDQuery, userID, limit, offset)
	if err != nil {
		r.log.Error(ctx, "failed to list offers by user", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	offers, err := scanOffersInFeed(rows)
	if err != nil {
		r.log.Error(ctx, "failed to scan offers for user feed", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}

	// Fetch total count for pagination metadata
	var total int
	err = r.db.QueryRowContext(ctx, countOffersByUserIDQuery, userID).Scan(&total)
	if err != nil {
		r.log.Warn(ctx, "failed to count total offers for user", zap.String("user_id", userID), zap.Error(err))
		total = len(offers) // fallback
	}

	return &domain.OffersInFeed{
		Meta: struct {
			Total  int
			Offset int
		}{
			Total:  total,
			Offset: offset,
		},
		Offers: offers,
	}, nil
}

// func (r *OfferRepository) listPhotosForOffer(ctx context.Context, offerID string) ([]string, error) {
// 	rows, err := r.db.QueryContext(ctx, listPhotosForOfferQuery, offerID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var urls []string
// 	for rows.Next() {
// 		var url string
// 		if err := rows.Scan(&url); err != nil {
// 			return nil, err
// 		}
// 		urls = append(urls, url)
// 	}
// 	return urls, nil
// }
