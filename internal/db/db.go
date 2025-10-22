package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func New(dataSourceName string) (*Database, error) {
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
	return &Database{db: db}, nil
}

func getMigrationsPath() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	candidates := []string{
		filepath.Join(wd, "internal", "db", "migrations"),
		filepath.Join(wd, "db", "migrations"),            
		filepath.Join(wd, "..", "internal", "db", "migrations"),     
		filepath.Join(wd, "..", "..", "internal", "db", "migrations"),
	}

	for _, path := range candidates {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			return "file://" + filepath.ToSlash(path)
		}
	}

	log.Fatalf("Migrations folder not found in any of the expected locations: %v", candidates)
	return ""
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

func (r *Database) GetDB() *sql.DB {
	return r.db
}

func (r *Database) MigrateDown(steps uint) error {
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
	
	return m.Steps(-int(steps))
}
