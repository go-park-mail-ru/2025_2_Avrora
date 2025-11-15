package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const (
	selectProfileWithUser = `
		p.id, p.user_id, p.first_name, p.last_name, p.phone, p.avatar_url,
		p.created_at, p.updated_at,
		u.email`

	getUserAndProfileLeftJoinQuery = `
		SELECT 
			p.id,
			u.id AS user_id,          
			p.first_name,
			p.last_name,
			p.phone,
			p.avatar_url,
			p.created_at,
			p.updated_at,
			u.email,
			u.role                  
		FROM users u
		LEFT JOIN profile p ON u.id = p.user_id
		WHERE u.id = $1`

	updateProfileQuery = `
		INSERT INTO profile (user_id, first_name, last_name, phone, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		ON CONFLICT (user_id) 
		DO UPDATE SET
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			phone = EXCLUDED.phone,
			avatar_url = EXCLUDED.avatar_url,
			updated_at = NOW();`

	updateUserPasswordHashQuery = `
		UPDATE users
		SET password_hash = $1, updated_at = $2
		WHERE id = $3`
)

type ProfileRepository struct {
	db  *pgxpool.Pool
	log *log.Logger
}

func NewProfileRepository(db *pgxpool.Pool, log *log.Logger) *ProfileRepository {
	return &ProfileRepository{db: db, log: log}
}

func (r *ProfileRepository) GetByUserID(ctx context.Context, userID string) (*domain.Profile, string, error) {
	var p domain.Profile
	var email, role string
	var id, userIDFromDB, firstName, lastName, phone, avatarURL *string
	var createdAt, updatedAt *time.Time

	err := r.db.QueryRow(ctx, getUserAndProfileLeftJoinQuery, userID).Scan(
		&id,
		&userIDFromDB,
		&firstName,
		&lastName,
		&phone,
		&avatarURL,
		&createdAt,
		&updatedAt,
		&email,
		&role,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", domain.ErrProfileNotFound // or domain.ErrUserNotFound?
		}
		r.log.Error(ctx, "failed to get user and profile", zap.String("user_id", userID), zap.Error(err))
		return nil, "", err
	}

	p = domain.Profile{
		ID:        SafeStringDeref(id),
		UserID:    SafeStringDeref(userIDFromDB),
		FirstName: SafeStringDeref(firstName),
		LastName:  SafeStringDeref(lastName),
		Phone:     SafeStringDeref(phone),
		AvatarURL: SafeStringDeref(avatarURL),
		Role:      role,
		Email:     email,
	}
	if createdAt != nil {
		p.CreatedAt = *createdAt
	}
	if updatedAt != nil {
		p.UpdatedAt = *updatedAt
	}

	return &p, email, nil
}

func (r *ProfileRepository) Update(ctx context.Context, userID string, upd *domain.ProfileUpdate) error {
	if upd == nil {
		return nil
	}

	_, err := r.db.Exec(ctx, updateProfileQuery,
		userID,
		upd.FirstName,
		upd.LastName,
		upd.Phone,
		upd.AvatarURL,
	)
	if err != nil {
		r.log.Error(ctx, "failed to upsert profile", zap.String("user_id", userID), zap.Error(err))
	}
	return err
}

func (r *ProfileRepository) UpdateSecurity(ctx context.Context, userID string, passwordHash string) error {
	now := time.Now().UTC()
	_, err := r.db.Exec(ctx, updateUserPasswordHashQuery,
		passwordHash,
		now,
		userID,
	)
	if err != nil {
		r.log.Error(ctx, "failed to update password", zap.String("user_id", userID), zap.Error(err))
	}
	return err
}

func (r *ProfileRepository) UpdateEmail(ctx context.Context, userID string, email string) error {
	_, err := r.db.Exec(ctx, updateUserEmailQuery, email, userID)
	if err != nil {
		r.log.Error(ctx, "failed to update email", zap.String("user_id", userID), zap.Error(err))
	}
	return err
}

func (r *ProfileRepository) GetUserByUserID(ctx context.Context, userID string) (*domain.User, error) {
	user := domain.User{}
	err := r.db.QueryRow(ctx, getUserByIDQuery, userID).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.Error(ctx, "failed to get user by ID", zap.String("id", userID), zap.Error(err))
			return nil, domain.ErrUserNotFound
		}
		r.log.Error(ctx, "failed to get user by ID", zap.String("id", userID), zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func SafeStringDeref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
