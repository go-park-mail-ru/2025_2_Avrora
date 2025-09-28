package db

import (
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestInitDB(t *testing.T) {
	t.Parallel()
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "avrora_test")
	// ... и т.д.
}

func TestDBConnection(t *testing.T) {
	dsn := "postgres://postgres:postgres@localhost/2025_2_Avrora_test?sslmode=disable"
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
	dsn := "postgres://postgres:postgres@localhost/2025_2_Avrora_test?sslmode=disable"

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
	dsn := "postgres://postgres:postgres@localhost/2025_2_Avrora_test?sslmode=disable"

	repo, err := New(dsn)
	if err != nil {
		t.Fatalf("Failed to initialize test DB: %v\n", err)
	}
	testDB := repo.GetDB()

	if testDB == nil {
		t.Fatal("could not getDB()")
	}
}
