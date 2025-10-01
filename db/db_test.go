package db

import (
	"os"
	"testing"

	"github.com/go-park-mail-ru/2025_2_Avrora/utils"
	_ "github.com/lib/pq"
)

// setupTestEnv выставляет переменные окружения для тестовой БД
func setupTestEnv() {
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "2025_2_Avrora_test")
}

// setupTestRepo инициализирует репозиторий и очищает все таблицы
func setupTestRepo(t *testing.T) *Repo {
	t.Helper()
	setupTestEnv()

	dsn := utils.GetPostgresDSN()
	t.Log("DSN:", dsn) // для проверки правильности DSN

	repo, err := New(dsn)
	if err != nil {
		t.Fatal("Failed to connect to test DB:", err)
	}

	if err := repo.ClearAllTables(); err != nil {
		t.Fatal("Failed to clear tables:", err)
	}

	return repo
}

func TestDBConnection(t *testing.T) {
	repo := setupTestRepo(t)
	defer repo.db.Close()

	if err := repo.db.Ping(); err != nil {
		t.Fatal("Ping failed:", err)
	}
}

func TestClearAllTables(t *testing.T) {
	repo := setupTestRepo(t)
	defer repo.db.Close()

	if err := repo.ClearAllTables(); err != nil {
		t.Fatal("Clear failed:", err)
	}
}

func TestGetDB(t *testing.T) {
	repo := setupTestRepo(t)
	defer repo.db.Close()

	if repo.GetDB() == nil {
		t.Fatal("could not getDB()")
	}
}
