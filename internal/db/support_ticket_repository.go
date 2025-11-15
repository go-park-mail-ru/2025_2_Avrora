package db

import (
	"context"
	"errors"
	"time"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// SQL constants for support tickets
const (
	createSupportTicketQuery = `
		INSERT INTO support_ticket (
			user_id,
			signed_email,
			response_email,
			name,
			category,
			description,
			status,
			created_at,
			updated_at
		) VALUES (
			$1,  -- user_id (can be NULL)
			$2,  -- signed_email
			$3,  -- response_email
			$4,  -- name
			$5,  -- category
			$6,  -- description
			$7,  -- status (default 'open')
			NOW(),
			NOW()
		)
		RETURNING id, created_at, updated_at
	`

	insertSupportTicketPhotoQuery = `
		INSERT INTO support_ticket_photo (ticket_id, photo_url, created_at)
		VALUES ($1, $2, NOW())
	`

	getSupportTicketByIDQuery = `
		SELECT
			st.id,
			st.user_id,
			st.signed_email,
			st.response_email,
			st.name,
			st.category,
			st.description,
			st.status,
			st.created_at,
			st.updated_at,
			COALESCE(
				ARRAY_AGG(stp.photo_url) FILTER (WHERE stp.photo_url IS NOT NULL),
				'{}'
			) AS photo_urls
		FROM support_ticket st
		LEFT JOIN support_ticket_photo stp ON stp.ticket_id = st.id
		WHERE st.id = $1
		GROUP BY st.id
	`

	getSupportTicketsByUserIDQuery = `
		SELECT
			st.id,
			st.user_id,
			st.signed_email,
			st.response_email,
			st.name,
			st.category,
			st.description,
			st.status,
			st.created_at,
			st.updated_at,
			COALESCE(
				ARRAY_AGG(stp.photo_url) FILTER (WHERE stp.photo_url IS NOT NULL),
				'{}'
			) AS photo_urls
		FROM support_ticket st
		LEFT JOIN support_ticket_photo stp ON stp.ticket_id = st.id
		WHERE st.user_id = $1
		GROUP BY st.id
		ORDER BY st.created_at DESC
		LIMIT $2 OFFSET $3
	`

	countSupportTicketsByUserIDQuery = `
		SELECT COUNT(*)
		FROM support_ticket
		WHERE user_id = $1
	`

	updateSupportTicketStatusQuery = `
		UPDATE support_ticket
		SET status = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING status, updated_at
	`

	deleteSupportTicketQuery = `
		DELETE FROM support_ticket
		WHERE id = $1
	`

	listAllSupportTicketsQuery = `
		SELECT
			st.id,
			st.user_id,
			st.signed_email,
			st.response_email,
			st.name,
			st.category,
			st.description,
			st.status,
			st.created_at,
			st.updated_at,
			COALESCE(
				ARRAY_AGG(stp.photo_url) FILTER (WHERE stp.photo_url IS NOT NULL),
				'{}'
			) AS photo_urls
		FROM support_ticket st
		LEFT JOIN support_ticket_photo stp ON stp.ticket_id = st.id
		GROUP BY st.id
		ORDER BY st.created_at DESC
		LIMIT $1 OFFSET $2
	`

	countAllSupportTicketsQuery = `
		SELECT COUNT(*)
		FROM support_ticket
	`
)

type SupportTicketRepository struct {
	db  *pgxpool.Pool
	log *log.Logger
}

func NewSupportTicketRepository(db *pgxpool.Pool, log *log.Logger) *SupportTicketRepository {
	return &SupportTicketRepository{db: db, log: log}
}

func scanSupportTicketRow(scanner interface {
	Scan(dest ...any) error
}) (*domain.SupportTicket, error) {
	var (
		userID      *string
		photoURLs   []string
		ticket      domain.SupportTicket
	)

	err := scanner.Scan(
		&ticket.ID,
		&userID,
		&ticket.SignedEmail,
		&ticket.ResponseEmail,
		&ticket.Name,
		&ticket.Category,
		&ticket.Description,
		&ticket.Status,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
		&photoURLs,
	)
	if err != nil {
		return nil, err
	}

	// Assign nullable fields
	ticket.UserID = userID
	ticket.PhotoURLs = photoURLs
	if ticket.PhotoURLs == nil {
		ticket.PhotoURLs = []string{}
	}

	return &ticket, nil
}

func scanSupportTicket(row pgx.Row) (*domain.SupportTicket, error) {
	return scanSupportTicketRow(row)
}

func scanSupportTickets(rows pgx.Rows) ([]domain.SupportTicket, error) {
	var tickets []domain.SupportTicket
	for rows.Next() {
		ticket, err := scanSupportTicketRow(rows)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, *ticket)
	}
	return tickets, rows.Err()
}

func (r *SupportTicketRepository) Create(ctx context.Context, ticket *domain.SupportTicket) error {
	now := time.Now().UTC()
	ticket.CreatedAt = now
	ticket.UpdatedAt = now

	tx, err := r.db.Begin(ctx)
	if err != nil {
		r.log.Error(ctx, "failed to begin transaction for create ticket", zap.Error(err))
		return err
	}
	defer tx.Rollback(ctx)

	// Handle nullable user_id
	var userIDParam any
	if ticket.UserID != nil {
		userIDParam = *ticket.UserID
	} else {
		userIDParam = nil
	}

	// Insert main ticket
	err = tx.QueryRow(ctx, createSupportTicketQuery,
		userIDParam,
		ticket.SignedEmail,
		ticket.ResponseEmail,
		ticket.Name,
		ticket.Category,
		ticket.Description,
		ticket.Status,
	).Scan(
		&ticket.ID,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
	)
	if err != nil {
		r.log.Error(ctx, "failed to create support ticket", zap.Error(err))
		return err
	}

	// Insert photos if any
	for _, url := range ticket.PhotoURLs {
		_, err := tx.Exec(ctx, insertSupportTicketPhotoQuery,
			ticket.ID,
			url,
		)
		if err != nil {
			r.log.Warn(ctx, "failed to insert ticket photo",
				zap.String("ticket_id", ticket.ID),
				zap.String("url", url),
				zap.Error(err))
			// Continue on photo error - ticket should still be created
		}
	}

	if err := tx.Commit(ctx); err != nil {
		r.log.Error(ctx, "failed to commit ticket creation transaction", zap.Error(err))
		return err
	}

	r.log.Info(ctx, "created support ticket", zap.String("id", ticket.ID))
	return nil
}

func (r *SupportTicketRepository) GetByID(ctx context.Context, id string) (*domain.SupportTicket, error) {
	ticket, err := scanSupportTicket(r.db.QueryRow(ctx, getSupportTicketByIDQuery, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSupportTicketNotFound
		}
		r.log.Error(ctx, "failed to get support ticket by ID", zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return ticket, nil
}

func (r *SupportTicketRepository) GetByUserID(ctx context.Context, userID string, page, limit int) ([]domain.SupportTicket, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	rows, err := r.db.Query(ctx, getSupportTicketsByUserIDQuery, userID, limit, offset)
	if err != nil {
		r.log.Error(ctx, "failed to get support tickets by user ID",
			zap.String("user_id", userID),
			zap.Int("page", page),
			zap.Int("limit", limit),
			zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	tickets, err := scanSupportTickets(rows)
	if err != nil {
		r.log.Error(ctx, "failed to scan support tickets", zap.Error(err))
		return nil, err
	}

	return tickets, nil
}

func (r *SupportTicketRepository) CountByUserID(ctx context.Context, userID string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, countSupportTicketsByUserIDQuery, userID).Scan(&count)
	if err != nil {
		r.log.Error(ctx, "failed to count support tickets by user ID",
			zap.String("user_id", userID),
			zap.Error(err))
		return 0, err
	}
	return count, nil
}

func (r *SupportTicketRepository) UpdateStatus(ctx context.Context, ticketID string, status domain.SupportTicketStatus) error {
	var (
		newStatus   string
		updatedAt   time.Time
	)

	err := r.db.QueryRow(ctx, updateSupportTicketStatusQuery,
		status,
		ticketID,
	).Scan(&newStatus, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrSupportTicketNotFound
		}
		r.log.Error(ctx, "failed to update support ticket status",
			zap.String("ticket_id", ticketID),
			zap.String("new_status", string(status)),
			zap.Error(err))
		return err
	}

	r.log.Info(ctx, "updated support ticket status",
		zap.String("ticket_id", ticketID),
		zap.String("new_status", newStatus))
	return nil
}

func (r *SupportTicketRepository) Delete(ctx context.Context, ticketID string) error {
	result, err := r.db.Exec(ctx, deleteSupportTicketQuery, ticketID)
	if err != nil {
		r.log.Error(ctx, "failed to delete support ticket",
			zap.String("ticket_id", ticketID),
			zap.Error(err))
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrSupportTicketNotFound
	}

	r.log.Info(ctx, "deleted support ticket",
		zap.String("ticket_id", ticketID),
		zap.Int64("rows_affected", rowsAffected))
	return nil
}

func (r *SupportTicketRepository) ListAll(ctx context.Context, page, limit int) ([]domain.SupportTicket, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	rows, err := r.db.Query(ctx, listAllSupportTicketsQuery, limit, offset)
	if err != nil {
		r.log.Error(ctx, "failed to list all support tickets",
			zap.Int("page", page),
			zap.Int("limit", limit),
			zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	tickets, err := scanSupportTickets(rows)
	if err != nil {
		r.log.Error(ctx, "failed to scan all support tickets", zap.Error(err))
		return nil, err
	}

	return tickets, nil
}

func (r *SupportTicketRepository) CountAll(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, countAllSupportTicketsQuery).Scan(&count)
	if err != nil {
		r.log.Error(ctx, "failed to count all support tickets", zap.Error(err))
		return 0, err
	}
	return count, nil
}