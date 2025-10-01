package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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

// buildMigrationPath формирует корректный file:// путь под текущую ОС
func buildMigrationPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("❌ Не удалось вычислить абсолютный путь для миграций: %v", err)
	}

	// Windows → заменяем слэши и добавляем "file:///"
	if runtime.GOOS == "windows" {
		abs = strings.ReplaceAll(abs, `\`, `/`)
		return fmt.Sprintf("file:///%s", abs)
	}

	// Linux/Mac → обычный путь
	return fmt.Sprintf("file://%s", abs)
}

func getMigrationsPath() string {
	// Берём текущий рабочий каталог
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Ищем папку migrations в нескольких стандартных местах
	candidates := []string{
		filepath.Join(wd, "db", "migrations"),             // основной вариант
		filepath.Join(wd, "..", "db", "migrations"),       // для тестов из /db
		filepath.Join(wd, "..", "..", "db", "migrations"), // на случай глубокого запуска
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

	// ⚡ убрал defer m.Close()

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

	// ⚡ убрал defer m.Close()
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
