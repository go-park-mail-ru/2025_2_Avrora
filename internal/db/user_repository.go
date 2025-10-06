package db

import (
	"database/sql"
	"errors"
	"time"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	user := domain.User{}
	err := r.db.QueryRow("SELECT id, email, password, created_at FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &domain.User{}, domain.ErrUserNotFound
		}
		return &domain.User{}, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *domain.User) error {
	now := time.Now()
	user.CreatedAt = now

	err := r.db.QueryRow(
		"INSERT INTO users (id, email, password, created_at) VALUES ($1, $2, $3, $4) RETURNING id",
		user.ID,
		user.Email,
		user.Password,
		user.CreatedAt,
	).Scan(&user.ID)

	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByID(id string) (*domain.User, error) {
	user := domain.User{}
	err := r.db.QueryRow("SELECT id, email, password, created_at FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &domain.User{}, domain.ErrUserNotFound
		}
		return &domain.User{}, err
	}
	return &user, nil
}