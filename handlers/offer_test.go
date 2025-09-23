package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/go-park-mail-ru/2025_2_Avrora/models"
)

func seedTestData(t *testing.T) {
	if err := testRepo.Init("postgres://postgres:postgres@localhost/2025_2_Avrora_test?sslmode=disable"); err != nil {
		t.Fatalf("Failed to reinitialize DB schema: %v", err)
	}

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

	if len(resp.Offers) != 1 {
		t.Errorf("Ожидался 1 ответ а получили %d", len(resp.Offers))
	}

	first := resp.Offers[0]
	if first.Title != "Продам квартиру на Тверской" {
		t.Errorf("не то название: %s", first.Title)
	}
	if first.Price != 100000000 {
		t.Errorf("не та цена: %d", first.Price)
	}
	if resp.Meta.Page != 1 || resp.Meta.Limit != 10 {
		t.Errorf("Expected page=1, limit=10; got page=%d, limit=%d", resp.Meta.Page, resp.Meta.Limit)
	}
}

func TestGetOffersHandler_Empty(t *testing.T) {
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