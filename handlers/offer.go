package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Avrora/db"
	"github.com/go-park-mail-ru/2025_2_Avrora/response"
)

func GetOffersHandler(w http.ResponseWriter, r *http.Request, repo *db.Repo) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, response.NewErrorResp("метод не поддерживается"))
		return
	}

	// Пагинация
	page, err := parseIntQueryParam(r, "page", 1)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := parseIntQueryParam(r, "limit", 20)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	offers, err := repo.Offer().FindAll(page, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, response.NewErrorResp("ошибка получения офферов"))
		return
	}

	total, err := repo.Offer().CountAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, response.NewErrorResp("считать не умею офферы :3"))
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"offers": offers,
		"meta": map[string]int{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func parseIntQueryParam(r *http.Request, key string, defaultValue int) (int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(param)
}