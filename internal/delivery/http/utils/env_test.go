package utils

import (
	"os"
	"path/filepath"
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
func TestFindProjectRoot_Found(t *testing.T) {
	tmpDir := t.TempDir()

	// создаём go.mod
	modPath := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(modPath, []byte("module test"), 0644); err != nil {
		t.Fatal(err)
	}

	// вложенная директория
	subDir := filepath.Join(tmpDir, "sub")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	root := findProjectRoot(subDir)
	if root != tmpDir {
		t.Errorf("expected %q, got %q", tmpDir, root)
	}
}
func TestValidateEnvVars_AllSet(t *testing.T) {
	for _, key := range RequiredEnvVars {
		os.Setenv(key, "value")
	}

	// Верю, что не должно упасть
	validateEnvVars()
}
