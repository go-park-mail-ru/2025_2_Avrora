package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Avrora/db"
	"github.com/go-park-mail-ru/2025_2_Avrora/response"
)

func GetOffersHandler(w http.ResponseWriter, r *http.Request, repo *db.Repo) {
	// Пагинация
	page, err := parseIntQueryParam(r, "page", 1)
	if err != nil || page < 1 {
		response.HandleError(w, err, http.StatusBadRequest, "невалидные параметры(page)")
	}

	limit, err := parseIntQueryParam(r, "limit", 20)
	if err != nil || limit < 1 || limit > 100 {
		response.HandleError(w, err, http.StatusBadRequest, "невалидные параметры(limit)")
	}

	offers, err := repo.Offer().FindAll(page, limit)
	if err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка бд при поиске офферов")
		return
	}

	total, err := repo.Offer().CountAll()
	if err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка бд при подсчете офферов")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{
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