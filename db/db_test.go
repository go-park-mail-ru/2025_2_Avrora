package db

import (
	"database/sql"
	"github.com/go-park-mail-ru/2025_2_Avrora/utils"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestInitDB(t *testing.T) {
	utils.LoadEnv() // <-- обязательно
	dsn := utils.GetPostgresDSN()
	t.Log("DB_USER:", os.Getenv("DB_USER"))
	t.Log("DB_PASSWORD:", os.Getenv("DB_PASSWORD"))
	t.Log("DB_NAME:", os.Getenv("DB_NAME"))
	repo, err := New(dsn)
	if err != nil {
		t.Fatal(err.Error())
	}

	if repo == nil {
		t.Fatal("testRepo is nil")
	}

	err = repo.GetDB().Ping()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestDBConnection(t *testing.T) {
	utils.LoadEnv()
	dsn := utils.GetPostgresDSN()
	repo, err := New(dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer repo.db.Close()

	err = repo.db.Ping()
	if err != nil {
		t.Fatal("Ping failed:", err)
	}

	err = repo.ClearAllTables()
	if err != nil {
		t.Fatal("Clear failed:", err)
	}
}

func TestClearAllTables(t *testing.T) {
	utils.LoadEnv()
	dsn := utils.GetPostgresDSN()

	repo, err := New(dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer repo.db.Close()

	testDB := repo.GetDB()

	err = testDB.Ping()
	if err != nil {
		t.Fatal(err)
	}

	err = repo.ClearAllTables()
	if err != nil {
		t.Fatal(err)
	}

}

func TestGetDB(t *testing.T) {
	utils.LoadEnv() // <- обязательно
	dsn := utils.GetPostgresDSN()

	repo, err := New(dsn)
	if err != nil {
		t.Fatalf("Failed to initialize test DB: %v\n", err)
	}
	testDB := repo.GetDB()

	if testDB == nil {
		t.Fatal("could not getDB()")
	}
}
func TestNew_InvalidDSN(t *testing.T) {
	_, err := New("user=wrong password=wrong dbname=wrong host=127.0.0.1 sslmode=disable")
	if err == nil {
		t.Fatal("expected error for invalid DSN, got nil")
	}
}
func TestClearAllTables_Empty(t *testing.T) {
	utils.LoadEnv()
	repo, _ := New(utils.GetPostgresDSN())
	defer repo.db.Close()

	repo.ClearAllTables()

	var count int
	err := repo.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal("users table not empty after ClearAllTables")
	}

	err = repo.db.QueryRow("SELECT COUNT(*) FROM offer").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal("offer table not empty after ClearAllTables")
	}
}
func TestApplyMigrations_ErrorIfDBClosed(t *testing.T) {
	utils.LoadEnv()
	repo, _ := New(utils.GetPostgresDSN())
	repo.db.Close() // имитируем закрытую БД

	err := applyMigrations(repo.db)
	if err == nil {
		t.Fatal("expected error for closed DB, got nil")
	}
}
func TestMigrateDown(t *testing.T) {
	utils.LoadEnv()
	repo, _ := New(utils.GetPostgresDSN())
	defer repo.db.Close()

	err := repo.MigrateDown(1)
	if err != nil {
		t.Fatal(err)
	}
}
func TestGetMigrationsPath_Exists(t *testing.T) {
	path := getMigrationsPath()
	if path == "" || path[:7] != "file://" {
		t.Fatal("invalid migrations path")
	}
}
func TestApplyMigrations_NoChange(t *testing.T) {
	utils.LoadEnv()

	// Применяем миграции дважды - второй раз должен вернуть ErrNoChange
	db, err := sql.Open("postgres", utils.GetPostgresDSN())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Первое применение
	err = applyMigrations(db)
	if err != nil {
		t.Fatal(err)
	}

	// Второе применение - должно вернуть nil (ErrNoChange обрабатывается)
	err = applyMigrations(db)
	if err != nil {
		t.Fatalf("Expected no error on second migration, got: %v", err)
	}
}
