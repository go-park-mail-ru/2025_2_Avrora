package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/go-park-mail-ru/2025_2_Avrora/db"
	"github.com/go-park-mail-ru/2025_2_Avrora/models"
	"github.com/go-park-mail-ru/2025_2_Avrora/utils"
)

func seedTestData(t *testing.T) {
	utils.LoadEnv()
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME_TEST"),
	)
	testRepo, _ := db.New(dsn)
	testRepo.ClearAllTables()
	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Join(filepath.Dir(filename), "..")
	sqlPath := filepath.Join(baseDir, "db", "seed", "mocks.sql")

	content, err := os.ReadFile(sqlPath)
	if err != nil {
		t.Fatalf("Failed to read mocks.sql: %v", err)
	}

	_, err = testRepo.GetDB().Exec(string(content))
	if err != nil {
		t.Fatalf("Failed to execute mocks.sql: %v", err)
	}

	t.Log("✅ Mock data seeded successfully")
}

func TestGetOffersHandler_Success(t *testing.T) {
	utils.LoadEnv()
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME_TEST"),
	)
	testRepo, _ := db.New(dsn)
	testRepo.ClearAllTables()
	seedTestData(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/offers?page=1&limit=10", nil)
	w := httptest.NewRecorder()

	GetOffersHandler(w, req, testRepo)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp struct {
		Offers []models.Offer `json:"offers"`
		Meta   struct {
			Page  int `json:"page"`
			Limit int `json:"limit"`
			Total int `json:"total"`
		} `json:"meta"`
	}

	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	if len(resp.Offers) != 5 {
		t.Errorf("Ожидался 5шт	 а получили %d", len(resp.Offers))
	}

	first := resp.Offers[0]
	if first.Title != "Уютная 2-комнатная квартира в центре" {
		t.Errorf("не то название: %s", first.Title)
	}
	if first.Price != 8500000 {
		t.Errorf("не та цена: %d", first.Price)
	}
	if resp.Meta.Page != 1 || resp.Meta.Limit != 10 {
		t.Errorf("Expected page=1, limit=10; got page=%d, limit=%d", resp.Meta.Page, resp.Meta.Limit)
	}
}

func TestGetOffersHandler_Empty(t *testing.T) {
	utils.LoadEnv()
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME_TEST"),
	)
	testRepo, _ := db.New(dsn)
	testRepo.ClearAllTables()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/offers", nil)
	w := httptest.NewRecorder()

	GetOffersHandler(w, req, testRepo)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp struct {
		Offers []models.Offer `json:"offers"`
		Meta   struct {
			Page  int `json:"page"`
			Limit int `json:"limit"`
			Total int `json:"total"`
		} `json:"meta"`
	}

	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	if len(resp.Offers) != 0 {
		t.Errorf("Expected 0 offers, got %d", len(resp.Offers))
	}
}