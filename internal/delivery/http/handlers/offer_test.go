package handlers

import (
	"bytes"
	"encoding/json"
	_ "errors"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// --- helper для создания offerHandler
func newTestOfferHandler() (*offerHandler, *mockOfferUsecase) {
	mockUC := new(mockOfferUsecase)
	logger := log.New(zap.NewNop())
	return &offerHandler{offerUsecase: mockUC, logger: logger}, mockUC
}

// --- TESTS ---

func TestGetOffers_Success(t *testing.T) {
	h, mockUC := newTestOfferHandler()
	expected := &domain.OffersInFeed{}
	mockUC.On("ListOffersInFeed", mock.Anything, 1, 10).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/offers?page=1&limit=10", nil)
	w := httptest.NewRecorder()

	h.GetOffers(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetOffers_BadPage(t *testing.T) {
	h, _ := newTestOfferHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/offers?page=abc", nil)
	w := httptest.NewRecorder()

	h.GetOffers(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateOffer_Success(t *testing.T) {
	h, mockUC := newTestOfferHandler()

	reqBody := CreateOfferRequest{
		Title: "Test", Description: "Desc", Address: "Addr", Price: 1000,
		Area: 50, Rooms: 2, OfferType: "sale", PropertyType: "flat", Floor: 3, TotalFloors: 5, UserID: "uid",
	}
	body, _ := json.Marshal(reqBody)
	mockUC.On("Create", mock.Anything, mock.AnythingOfType("*domain.Offer")).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/offers/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateOffer(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	mockUC.AssertExpectations(t)
}

func TestCreateOffer_InvalidJSON(t *testing.T) {
	h, _ := newTestOfferHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/offers/create", bytes.NewReader([]byte("{bad json")))
	w := httptest.NewRecorder()
	h.CreateOffer(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateOffer_InvalidInput(t *testing.T) {
	h, mockUC := newTestOfferHandler()
	reqBody := CreateOfferRequest{Title: "T", Address: "A"}
	body, _ := json.Marshal(reqBody)
	mockUC.On("Create", mock.Anything, mock.AnythingOfType("*domain.Offer")).Return(usecase.ErrInvalidInput)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/offers/create", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateOffer(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteOffer_Success(t *testing.T) {
	h, mockUC := newTestOfferHandler()
	mockUC.On("Delete", mock.Anything, "123").Return(nil)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/offers/delete/123", nil)
	w := httptest.NewRecorder()
	h.DeleteOffer(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestDeleteOffer_NoID(t *testing.T) {
	h, _ := newTestOfferHandler()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/offers/delete/", nil)
	w := httptest.NewRecorder()
	h.DeleteOffer(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateOffer_Success(t *testing.T) {
	h, mockUC := newTestOfferHandler()
	reqBody := UpdateOfferRequest{
		Title: "Upd", Address: "Addr", Price: 2000, Area: 60, Rooms: 3,
		Floor: 2, TotalFloors: 5, UserID: "uid",
	}
	body, _ := json.Marshal(reqBody)
	mockUC.On("Update", mock.Anything, mock.AnythingOfType("*domain.Offer")).Return(nil)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/offers/update/abc", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.UpdateOffer(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateOffer_InvalidJSON(t *testing.T) {
	h, _ := newTestOfferHandler()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/offers/update/abc", bytes.NewReader([]byte("{bad json")))
	w := httptest.NewRecorder()
	h.UpdateOffer(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateOffer_NoID(t *testing.T) {
	h, _ := newTestOfferHandler()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/offers/update/", bytes.NewReader([]byte("{}")))
	w := httptest.NewRecorder()
	h.UpdateOffer(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetOffer_Success(t *testing.T) {
	h, mockUC := newTestOfferHandler()
	expected := &domain.Offer{ID: "123"}
	mockUC.On("Get", mock.Anything, "123").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/offers/123", nil)
	w := httptest.NewRecorder()
	h.GetOffer(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetOffer_NoID(t *testing.T) {
	h, _ := newTestOfferHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/offers", nil)
	w := httptest.NewRecorder()
	h.GetOffer(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetMyOffers_Success(t *testing.T) {
	h, mockUC := newTestOfferHandler()
	expected := &domain.OffersInFeed{}
	mockUC.On("ListOffersInFeedByUserID", mock.Anything, "user123", 1, 10).Return(expected, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/profile/myoffers/user123", nil)
	w := httptest.NewRecorder()
	h.GetMyOffers(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetMyOffers_NoUserID(t *testing.T) {
	h, _ := newTestOfferHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/profile/myoffers/", nil)
	w := httptest.NewRecorder()
	h.GetMyOffers(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
