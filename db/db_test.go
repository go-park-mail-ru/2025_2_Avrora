package db

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-park-mail-ru/2025_2_Avrora/utils"
	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	utils.LoadEnv()
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)
	r := NewRepo()
	r.Init(dsn)
	DB := r.db

	code := m.Run()

	_, _ = DB.Exec("DROP TABLE IF EXISTS users")
	os.Exit(code)
}

func TestInitDB(t *testing.T) {
	dsn := "postgres://postgres:postgres@localhost/2025_2_Avrora_test?sslmode=disable"
	r := NewRepo()
	r.Init(dsn)
}

func TestClearAllTables(t *testing.T) {
	dsn := "postgres://postgres:postgres@localhost/2025_2_Avrora_test?sslmode=disable"
	r := NewRepo()
	r.Init(dsn)
	r.ClearAllTables()
}

func TestCreateTables(t *testing.T) {
	dsn := "postgres://postgres:postgres@localhost/2025_2_Avrora_test?sslmode=disable"
	r := NewRepo()
	r.Init(dsn)
	r.createTables()
}

func TestGetDB(t *testing.T) {
	dsn := "postgres://postgres:postgres@localhost/2025_2_Avrora_test?sslmode=disable"
	r := NewRepo()
	r.Init(dsn)
	r.GetDB()
}