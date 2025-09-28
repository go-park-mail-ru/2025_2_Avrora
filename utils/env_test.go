package utils

import (
	"os"
	"testing"
)

func TestGetPostgresDSN(t *testing.T) {
	os.Setenv("DB_USER", "user")
	os.Setenv("DB_PASSWORD", "pass")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "mydb")

	expected := "postgres://user:pass@localhost:5432/mydb?sslmode=disable"
	dsn := GetPostgresDSN()
	if dsn != expected {
		t.Errorf("expected %q, got %q", expected, dsn)
	}
}
