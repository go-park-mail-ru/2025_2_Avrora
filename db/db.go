package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	_ "github.com/lib/pq"
)

type Repo struct {
	db *sql.DB
	once sync.Once
}

func NewRepo() *Repo {
	return &Repo{}
}

func (r *Repo) Init(dataSourceName string) error {
	var err error
	r.once.Do(func() {
		r.db, err = sql.Open("postgres", dataSourceName)
		if err != nil {
			return
		}

		err = r.createTables()
		if err != nil {
			return
		}

		err = r.db.Ping()
		if err != nil {
			return
		}

		log.Println("Инициализация базы данных завершена")
	})
	return err
}

func getMigrationsPath() string {
	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Join(filepath.Dir(filename), "..")
	return filepath.Join(baseDir, "db", "migrations")
}

func (r *Repo) createTables() error {
	migrationsPath := getMigrationsPath()
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") && len(file.Name()) >= 8 {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}

	sort.Strings(sqlFiles)

	log.Println("Применяем миграции...")
	for _, filename := range sqlFiles {
		fullPath := filepath.Join(migrationsPath, filename)

		content, err := os.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		_, err = r.db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}

		log.Printf("✅ Выполнена миграция %s", filename)
	}

	log.Println("✅ Все миграции применены")
	return nil
}

func (r *Repo) GetDB() *sql.DB {
	return r.db
}

func (r *Repo) ClearAllTables() {
	_, _ = r.GetDB().Exec("DELETE FROM photo")
	_, _ = r.GetDB().Exec("DELETE FROM offer")
	_, _ = r.GetDB().Exec("DELETE FROM location")
	_, _ = r.GetDB().Exec("DELETE FROM region")
	_, _ = r.GetDB().Exec("DELETE FROM category")
	_, _ = r.GetDB().Exec("DELETE FROM user")
}