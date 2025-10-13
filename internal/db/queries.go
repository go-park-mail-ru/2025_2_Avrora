package db

// User queries

const (
	getUserByEmailQuery = "SELECT id, email, password, created_at FROM users WHERE email = $1"
	getUserByIDQuery    = "SELECT id, email, password, created_at FROM users WHERE id = $1"
	createUserQuery     = "INSERT INTO users (email, password, created_at) VALUES ($1, $2, $3) RETURNING id"
)

// Offer queries

const (
	getOfferByIDQuery = `
		SELECT id, user_id, location_id, category_id, title, description, image, price, area, rooms, address, offer_type, created_at, updated_at
		FROM offer
		WHERE id = $1`

	listOffersQuery = `
		SELECT id, user_id, location_id, category_id, title, description, image, price, area, rooms, address, offer_type, created_at, updated_at
		FROM offer
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	createOfferQuery = `
		INSERT INTO offer (id, user_id, location_id, category_id, title, description, image, price, area, rooms, address, offer_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id`

	updateOfferQuery = `
		UPDATE offer
		SET title = $1, description = $2, image = $3, price = $4, area = $5, rooms = $6, address = $7, offer_type = $8, updated_at = $9
		WHERE id = $10`

	deleteOfferQuery = "DELETE FROM offer WHERE id = $1"

	countAllOffersQuery = "SELECT COUNT(*) FROM offer"

	listOffersByUserIDQuery = `
		SELECT id, user_id, location_id, category_id, title, description, image, price, area, rooms, address, offer_type, created_at, updated_at
		FROM offer
		WHERE user_id = $1
		ORDER BY created_at DESC`
)