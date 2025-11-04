package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

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
	db  *pgxpool.Pool
	Log *log.Logger
}

func NewUserRepository(db *pgxpool.Pool, log *log.Logger) *UserRepository {
	return &UserRepository{db: db, Log: log}
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := r.db.QueryRow(ctx, getUserByEmailQuery, email)
	user := domain.User{}
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		r.Log.Error(ctx, "failed to get user by email", zap.String("email", email), zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now

	var id string
	err := r.db.QueryRow(ctx, createUserQuery, user.Email, user.PasswordHash, now).Scan(&id)
	if err != nil {
		r.Log.Error(ctx, "failed to create user", zap.Error(err))
		return err
	}
	user.ID = id
	return nil
}

func (r *UserRepository) UpdateEmail(ctx context.Context, id, email string) error {
	_, err := r.db.Exec(ctx, updateUserEmailQuery, email, id)
	if err != nil {
		r.Log.Error(ctx, "failed to update email", zap.String("id", id), zap.Error(err))
	}
	return err
}
