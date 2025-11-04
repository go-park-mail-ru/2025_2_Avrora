package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// ---- Тест успешной фильтрации ----
func TestOfferHandler_FilterOffers_Success(t *testing.T) {
	mockUC := new(mockOfferUsecase)
	zapLogger, _ := zap.NewDevelopment()
	logger := log.New(zapLogger)
	h := NewOfferHandler(mockUC, logger)

	expectedOffers := []domain.OfferInFeed{
		{
			ID:       "1",
			UserID:   "u1",
			Address:  "Test Street",
			Price:    12345,
			Rooms:    2,
			Area:     42.5,
			ImageURL: "img.jpg",
		},
	}

	mockUC.On("FilterOffers", mock.Anything, mock.AnythingOfType("*domain.OfferFilter"), 20, 0).
		Return(expectedOffers, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/offers?address=Test+Street", nil)
	w := httptest.NewRecorder()

	h.FilterOffers(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var got []domain.OfferInFeed
	err := json.NewDecoder(res.Body).Decode(&got)
	assert.NoError(t, err)
	assert.Equal(t, expectedOffers, got)

	mockUC.AssertExpectations(t)
}

// ---- Тест ошибки фильтрации ----
func TestOfferHandler_FilterOffers_Error(t *testing.T) {
	mockUC := new(mockOfferUsecase)
	zapLogger, _ := zap.NewDevelopment()
	logger := log.New(zapLogger)
	h := NewOfferHandler(mockUC, logger)

	mockUC.On("FilterOffers", mock.Anything, mock.AnythingOfType("*domain.OfferFilter"), 20, 0).
		Return(nil, errors.New("db error")).Once()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/offers?address=ErrorStreet", nil)
	w := httptest.NewRecorder()

	h.FilterOffers(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	mockUC.AssertExpectations(t)
}
