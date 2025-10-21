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

type UserRepository struct {
	db *sql.DB
	log *log.Logger
}

func NewUserRepository(db *sql.DB, log *log.Logger) *UserRepository {
	return &UserRepository{db: db, log: log}
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := domain.User{}
	err := r.db.QueryRow("SELECT id, email, password, created_at FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.Error(ctx, "user not found", zap.String("email", email))
			return &domain.User{}, domain.ErrUserNotFound
		}
		r.log.Error(ctx, "failed to get user", zap.String("email", email), zap.Error(err))
		return &domain.User{}, err
	}

	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	now := time.Now()
	user.CreatedAt = now

	err := r.db.QueryRow(
		"INSERT INTO users (email, password, created_at) VALUES ($1, $2, $3) RETURNING id",
		user.Email,
		user.Password,
		user.CreatedAt,
	).Scan(&user.ID)

	if err != nil {
		r.log.Error(ctx, "failed to create user", zap.Error(err))
		return err
	}

	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user := domain.User{}
	err := r.db.QueryRow("SELECT id, email, password, created_at FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.Error(ctx, "user not found", zap.String("id", id))
			return &domain.User{}, domain.ErrUserNotFound
		}
		r.log.Error(ctx, "failed to get user", zap.String("id", id), zap.Error(err))
		return &domain.User{}, err
	}
	
	return &user, nil
}