package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_Avrora/models"
)

func TestGetOffersHandler(t *testing.T) {
	// БД не используется — только моки
	req := httptest.NewRequest(http.MethodGet, "/offers?page=1&limit=2", nil)
	w := httptest.NewRecorder()

	GetOffersHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp struct {
		Offers []models.Offer `json:"offers"`
		Page   int            `json:"page"`
		Limit  int            `json:"limit"`
		Total  int            `json:"total"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	if len(resp.Offers) != 2 {
		t.Errorf("Expected 2 offers, got %d", len(resp.Offers))
	}

	if resp.Total != 3 {
		t.Errorf("Expected total 3, got %d", resp.Total)
	}

	first := resp.Offers[0]
	if first.Title != "Продам 2-комнатную квартиру в центре" {
		t.Errorf("Unexpected offer title: %s", first.Title)
	}
}
