package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"database/sql"

	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Database struct {
	db *pgxpool.Pool
}

func New(dataSourceName string) (*Database, error) {
	if err := applyMigrations(dataSourceName); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	pool, err := pgxpool.New(context.Background(), dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	log.Println("✅ Инициализация базы данных завершена")
	return &Database{db: pool}, nil
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

func applyMigrations(dataSourceName string) error {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return fmt.Errorf("failed to open db for migrations: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping db for migrations: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		getMigrationsPath(),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrations failed: %w", err)
	}

	log.Println("✅ Все миграции применены")
	return nil
}

func (d *Database) GetDB() *pgxpool.Pool {
	return d.db
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) MigrateDown(steps uint) error {
	db, err := sql.Open("pgx", d.db.Config().ConnConfig.Copy().ConnString())
	if err != nil {
		return err
	}
	defer db.Close()

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

	return m.Steps(-int(steps))
}