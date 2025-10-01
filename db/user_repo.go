package db

import (
	"database/sql"

	"github.com/go-park-mail-ru/2025_2_Avrora/models"
)

type UserRepo struct {
	db *sql.DB
}

func (r *Repo) User() *UserRepo {
	return &UserRepo{db: r.GetDB()}
}

func (ur *UserRepo) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := ur.db.QueryRow("SELECT id, email, password FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepo) Create(user *models.User) error {
	return ur.db.QueryRow(
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id",
		user.Email, user.Password,
	).Scan(&user.ID)
}

func (ur *UserRepo) ClearUserTable() error {
	_, err := ur.db.Exec("DELETE FROM users")
	return err
}
