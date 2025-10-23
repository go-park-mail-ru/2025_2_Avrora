package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"go.uber.org/zap"
)

const (
	getUserByEmailQuery = `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	getUserByIDQuery = `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	createUserQuery = `
		INSERT INTO users (email, password_hash, created_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	updateUserEmailQuery = `
		UPDATE users
		SET email = $1
		WHERE id = $2
	`
)

type UserRepository struct {
	db  *sql.DB
	log *log.Logger
}

func NewUserRepository(db *sql.DB, log *log.Logger) *UserRepository {
	return &UserRepository{db: db, log: log}
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := domain.User{}
	err := r.db.QueryRowContext(ctx, getUserByEmailQuery, email).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		r.log.Error(ctx, "failed to get user by email", zap.String("email", email), zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := r.db.QueryRowContext(
		ctx,
		createUserQuery,
		user.Email,
		user.PasswordHash,
		now,
	).Scan(&user.ID)

	if err != nil {
		r.log.Error(ctx, "failed to create user", zap.Error(err))
		return err
	}

	return nil
}

func (r* UserRepository) UpdateEmail(ctx context.Context, id, email string) error {
	_, err := r.db.ExecContext(ctx, updateUserEmailQuery, email, id)
	if err != nil {
		r.log.Error(ctx, "failed to update email", zap.String("id", id), zap.Error(err))
	}
	return err
}