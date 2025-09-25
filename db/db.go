package db

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Repo struct {
	db *sql.DB
}

func New(dataSourceName string) (*Repo, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := applyMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Инициализация базы данных завершена")
	return &Repo{db: db}, nil
}

func getMigrationsPath() string {
	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Join(filepath.Dir(filename), "..")
	return "file://" + filepath.Join(baseDir, "db", "migrations")
}

func applyMigrations(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return fmt.Errorf("DB closed before migrations: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		getMigrationsPath(),
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("DB closed after migrations: %w", err)
	}

	log.Println("✅ Все миграции применены")
	return nil
}

func (r *Repo) GetDB() *sql.DB {
	return r.db
}

func (r *Repo) MigrateDown(steps uint) error {
	driver, err := postgres.WithInstance(r.db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		getMigrationsPath(),
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}
	defer m.Close()

	return m.Steps(-int(steps))
}


func (r *Repo) ClearAllTables() error {
	_, err := r.db.Exec(`
		TRUNCATE TABLE offer, users
		RESTART IDENTITY
		CASCADE
	`)
	return err
}