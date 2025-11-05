package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	_ "errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// ---- Мок для IComplexUsecase ----

type mockComplexUsecase struct {
	mock.Mock
}

func (m *mockComplexUsecase) GetByID(ctx context.Context, id string) (*domain.HousingComplex, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.HousingComplex), args.Error(1)
}

func (m *mockComplexUsecase) List(ctx context.Context, page, limit int) (*domain.ComplexesInFeed, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ComplexesInFeed), args.Error(1)
}

func (m *mockComplexUsecase) Create(ctx context.Context, complex *domain.HousingComplex) error {
	args := m.Called(ctx, complex)
	return args.Error(0)
}

func (m *mockComplexUsecase) Update(ctx context.Context, complex *domain.HousingComplex) error {
	args := m.Called(ctx, complex)
	return args.Error(0)
}

func (m *mockComplexUsecase) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---- Вспомогательная функция ----

func newTestComplexHandler() *ComplexHandler {
	mockUC := new(mockComplexUsecase)
	logger := log.New(zap.NewNop())
	return &ComplexHandler{
		complexUsecase: mockUC,
		logger:         logger,
	}
}

// ---- Тесты ----

func TestGetComplexByID_Success(t *testing.T) {
	mockUC := new(mockComplexUsecase)
	logger := log.New(zap.NewNop())
	h := &ComplexHandler{complexUsecase: mockUC, logger: logger}

	id := uuid.NewString()
	expected := &domain.HousingComplex{ID: id, Name: "Test Complex"}

	mockUC.On("GetByID", mock.Anything, id).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/complexes/"+id, nil)
	w := httptest.NewRecorder()

	h.GetComplexByID(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp domain.HousingComplex
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, resp.ID)
	mockUC.AssertExpectations(t)
}

func TestGetComplexByID_NotFound(t *testing.T) {
	mockUC := new(mockComplexUsecase)
	logger := log.New(zap.NewNop())
	h := &ComplexHandler{complexUsecase: mockUC, logger: logger}

	id := uuid.NewString()
	mockUC.On("GetByID", mock.Anything, id).Return(nil, domain.ErrComplexNotFound)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/complexes/"+id, nil)
	w := httptest.NewRecorder()

	h.GetComplexByID(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUC.AssertExpectations(t)
}

func TestListComplexes_Success(t *testing.T) {
	mockUC := new(mockComplexUsecase)
	logger := log.New(zap.NewNop())
	h := &ComplexHandler{complexUsecase: mockUC, logger: logger}

	expected := &domain.ComplexesInFeed{}
	expected.Meta.Total = 1

	mockUC.On("List", mock.Anything, 1, 10).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/complexes?page=1&limit=10", nil)
	w := httptest.NewRecorder()

	h.ListComplexes(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp domain.ComplexesInFeed
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.Meta.Total)

	mockUC.AssertExpectations(t)
}

func TestCreateComplex_Success(t *testing.T) {
	mockUC := new(mockComplexUsecase)
	logger := log.New(zap.NewNop())
	h := &ComplexHandler{complexUsecase: mockUC, logger: logger}

	reqBody := CreateComplexRequest{
		Name: "New Complex",
	}
	body, _ := json.Marshal(reqBody)

	mockUC.On("Create", mock.Anything, mock.AnythingOfType("*domain.HousingComplex")).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/complexes", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateComplex(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockUC.AssertExpectations(t)
}

func TestDeleteComplex_Success(t *testing.T) {
	mockUC := new(mockComplexUsecase)
	logger := log.New(zap.NewNop())
	h := &ComplexHandler{complexUsecase: mockUC, logger: logger}

	id := uuid.NewString()
	mockUC.On("Delete", mock.Anything, id).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/complexes/delete/"+id, nil)
	w := httptest.NewRecorder()

	h.DeleteComplex(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockUC.AssertExpectations(t)
}

func TestDeleteComplex_NotFound(t *testing.T) {
	mockUC := new(mockComplexUsecase)
	logger := log.New(zap.NewNop())
	h := &ComplexHandler{complexUsecase: mockUC, logger: logger}

	id := uuid.NewString()
	mockUC.On("Delete", mock.Anything, id).Return(domain.ErrComplexNotFound)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/complexes/delete/"+id, nil)
	w := httptest.NewRecorder()

	h.DeleteComplex(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUC.AssertExpectations(t)
}
