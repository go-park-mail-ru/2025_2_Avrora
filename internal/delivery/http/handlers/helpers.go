package handlers

import (
	"net/http"
	"strconv"
)

type DeleteOfferRequest struct {
	ID int `json:"id"`
}

func parseIntQueryParam(r *http.Request, key string, defaultValue int) (int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(param)
}