package db

import (
	"testing"

	_ "github.com/lib/pq"
)

func TestInitDB(t *testing.T) {
	dsn := "postgres://postgres:postgres@localhost/2025_2_Avrora_test?sslmode=disable"

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
