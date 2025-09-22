package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2025_2_Avrora/models"
	"github.com/go-park-mail-ru/2025_2_Avrora/db"
)

var mockOffers = []models.Offer{
	{
		ID:          1,
		UserID:      1,
		Title:       "Продам 2-комнатную квартиру в центре",
		Description: "Свежий ремонт, вид на парк, развитая инфраструктура.",
		Price:       125000000,
		Area:        54.5,
		Rooms:       2,
		Address:     "Москва, Тверская улица, 15",
		OfferType:   "sale",
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now().Add(-12 * time.Hour),
	},
	{
		ID:          2,
		UserID:      2,
		Title:       "Сдам уютную студию у метро",
		Description: "Все включено, Wi-Fi, консьерж, парковка.",
		Price:       6500000, // 65 000 ₽/мес
		Area:        28.0,
		Rooms:       1,
		Address:     "Санкт-Петербург, Невский проспект, 22",
		OfferType:   "rent",
		CreatedAt:   time.Now().Add(-12 * time.Hour),
		UpdatedAt:   time.Now().Add(-6 * time.Hour),
	},
	{
		ID:          3,
		UserID:      3,
		Title:       "Квартира улучшенной планировки",
		Description: "Панорамные окна, охраняемая территория, детская площадка.",
		Price:       98000000,
		Area:        72.3,
		Rooms:       3,
		Address:     "Екатеринбург, ул. Малышева, 50",
		OfferType:   "sale",
		CreatedAt:   time.Now().Add(-48 * time.Hour),
		UpdatedAt:   time.Now().Add(-24 * time.Hour),
	},
}

const useMockOffers = true

func GetOffersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Пагинация
	page := 1
	limit := 20

	if p := r.URL.Query().Get("page"); p != "" {
		if pageNum, err := strconv.Atoi(p); err == nil && pageNum > 0 {
			page = pageNum
		}
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		if lim, err := strconv.Atoi(l); err == nil && lim > 0 && lim <= 100 {
			limit = lim
		}
	}

	offset := (page - 1) * limit

	var offers []models.Offer

	if useMockOffers {
		start := offset
		end := offset + limit
		if start > len(mockOffers) {
			offers = []models.Offer{}
		} else {
			if end > len(mockOffers) {
				end = len(mockOffers)
			}
			offers = mockOffers[start:end]
		}
	} else {
		rows, err := db.DB.Query(`
			SELECT id, user_id, title, description, price, area, rooms, address, offer_type, created_at, updated_at
			FROM offers
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`, limit, offset)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch offers"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var o models.Offer
			err := rows.Scan(
				&o.ID,
				&o.UserID,
				&o.Title,
				&o.Description,
				&o.Price,
				&o.Area,
				&o.Rooms,
				&o.Address,
				&o.OfferType,
				&o.CreatedAt,
				&o.UpdatedAt,
			)
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to scan offer"})
				return
			}
			offers = append(offers, o)
		}

		if err = rows.Err(); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Row iteration error"})
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]any {
		"offers": offers,
		"page":   page,
		"limit":  limit,
		"total":  len(mockOffers), // в реальной БД нужно COUNT(*)
	})
}