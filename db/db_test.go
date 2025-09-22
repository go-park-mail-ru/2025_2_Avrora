package db

import (
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	dsn := "postgres://postgres:postgres@localhost/2025_2_Avrora?sslmode=disable"
	InitDB(dsn)

	code := m.Run()

	_, _ = DB.Exec("DROP TABLE IF EXISTS users")
	os.Exit(code)
}

func TestInitDB(t *testing.T) {
	dsn := "postgres://postgres:postgres@localhost/2025_2_Avrora?sslmode=disable"
	InitDB(dsn)
}