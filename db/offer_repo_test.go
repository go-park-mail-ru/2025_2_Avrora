package db

import (
	"testing"

	"github.com/go-park-mail-ru/2025_2_Avrora/models"
)

// createTestUser создаёт пользователя для внешнего ключа в offer
func createTestUser(t *testing.T, repo *Repo) int {
	var id int
	err := repo.GetDB().QueryRow(
		`INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`,
		"user@test.com", "password123",
	).Scan(&id)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	return id
}

func TestOfferRepo_CreateAndFindByID(t *testing.T) {
	repo := setupTestRepo(t)
	defer repo.db.Close()

	userID := createTestUser(t, repo)
	or := repo.Offer()

	offer := &models.Offer{
		UserID:      userID,
		LocationID:  1,
		CategoryID:  1,
		Title:       "Test Offer",
		Description: "Some description",
		Price:       1000,
		Area:        55.5,
		Rooms:       3,
		Address:     "Test street 1",
		OfferType:   "sale",
	}

	if err := or.Create(offer); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if offer.ID == 0 {
		t.Fatal("expected non-zero ID after Create")
	}

	found, err := or.FindByID(offer.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found == nil || found.Title != offer.Title {
		t.Fatalf("unexpected found offer: %+v", found)
	}
}

func TestOfferRepo_FindAll(t *testing.T) {
	repo := setupTestRepo(t)
	defer repo.db.Close()

	userID := createTestUser(t, repo)
	or := repo.Offer()

	// создаём несколько офферов
	for i := 0; i < 3; i++ {
		_ = or.Create(&models.Offer{
			UserID:     userID,
			LocationID: 1,
			CategoryID: 1,
			Title:      "Offer",
			Price:      100 * (i + 1),
			Area:       20,
			Rooms:      i,
			Address:    "Addr",
			OfferType:  "rent",
		})
	}

	offers, err := or.FindAll(1, 10)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	if len(offers) != 3 {
		t.Fatalf("expected 3 offers, got %d", len(offers))
	}
}

func TestOfferRepo_Update(t *testing.T) {
	repo := setupTestRepo(t)
	defer repo.db.Close()

	userID := createTestUser(t, repo)
	or := repo.Offer()

	offer := &models.Offer{
		UserID:     userID,
		LocationID: 1,
		CategoryID: 1,
		Title:      "Old title",
		Price:      500,
		Area:       33,
		Rooms:      1,
		Address:    "Addr",
		OfferType:  "rent",
	}
	_ = or.Create(offer)

	offer.Title = "New title"
	offer.Price = 999

	if err := or.Update(offer); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated, _ := or.FindByID(offer.ID)
	if updated.Title != "New title" || updated.Price != 999 {
		t.Fatalf("offer not updated: %+v", updated)
	}
}

func TestOfferRepo_CountAllAndClear(t *testing.T) {
	repo := setupTestRepo(t)
	defer repo.db.Close()

	userID := createTestUser(t, repo)
	or := repo.Offer()

	for i := 0; i < 2; i++ {
		_ = or.Create(&models.Offer{
			UserID:     userID,
			LocationID: 1,
			CategoryID: 1,
			Title:      "Offer",
			Price:      100,
			Area:       40,
			Rooms:      2,
			Address:    "Addr",
			OfferType:  "sale",
		})
	}

	total, err := or.CountAll()
	if err != nil {
		t.Fatalf("CountAll failed: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected 2 offers, got %d", total)
	}

	or.ClearOfferTable()
	total, _ = or.CountAll()
	if total != 0 {
		t.Fatalf("expected 0 after clear, got %d", total)
	}
}
