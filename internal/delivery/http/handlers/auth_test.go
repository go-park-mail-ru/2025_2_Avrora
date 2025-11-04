package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/response"
	mylog "github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	usecase "github.com/go-park-mail-ru/2025_2_Avrora/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// -------------------- Моки --------------------

type mockAuthUsecase struct {
	mock.Mock
}

func (m *mockAuthUsecase) Register(ctx context.Context, email, password string) error {
	args := m.Called(ctx, email, password)
	return args.Error(0)
}

func (m *mockAuthUsecase) Login(ctx context.Context, email, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}

func (m *mockAuthUsecase) Logout(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

// -------------------- Вспомогательные --------------------

func newTestHandler() *authHandler {
	mockUC := new(mockAuthUsecase)
	zapLogger := zap.NewNop()
	appLogger := mylog.New(zapLogger)
	return &authHandler{
		authUsecase: mockUC,
		logger:      appLogger, // ✅ тип совпадает (*log.Logger из твоего пакета)
	}
}

// -------------------- Тесты Register --------------------

func TestRegister_Success(t *testing.T) {
	h := newTestHandler()
	mockUC := h.authUsecase.(*mockAuthUsecase)

	reqBody := RegisterRequest{
		Email:    "test@example.com",
		Password: "123456A$aa",
	}
	body, _ := json.Marshal(reqBody)

	mockUC.On("Register", mock.Anything, reqBody.Email, reqBody.Password).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Register(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp RegisterResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, reqBody.Email, resp.Email)
	mockUC.AssertExpectations(t)
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	h := newTestHandler()
	mockUC := h.authUsecase.(*mockAuthUsecase)

	reqBody := RegisterRequest{
		Email:    "exists@example.com",
		Password: "123456A$aa",
	}
	body, _ := json.Marshal(reqBody)

	mockUC.On("Register", mock.Anything, reqBody.Email, reqBody.Password).
		Return(usecase.ErrUserAlreadyExists)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Register(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	mockUC.AssertExpectations(t)
}

// -------------------- Тесты Login --------------------

func TestLogin_Success(t *testing.T) {
	h := newTestHandler()
	mockUC := h.authUsecase.(*mockAuthUsecase)

	reqBody := LoginRequest{
		Email:    "test@example.com",
		Password: "123456A$aa",
	}
	body, _ := json.Marshal(reqBody)

	mockUC.On("Login", mock.Anything, reqBody.Email, reqBody.Password).
		Return("mock_token", nil)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Login(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp AuthResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "mock_token", resp.Token)
	mockUC.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	h := newTestHandler()
	mockUC := h.authUsecase.(*mockAuthUsecase)

	reqBody := LoginRequest{
		Email:    "bad@example.com",
		Password: "123456A$aa2",
	}
	body, _ := json.Marshal(reqBody)

	mockUC.On("Login", mock.Anything, reqBody.Email, reqBody.Password).
		Return("", usecase.ErrInvalidCredentials)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Login(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockUC.AssertExpectations(t)
}

// -------------------- Тесты Logout --------------------

func TestLogout_Success(t *testing.T) {
	h := newTestHandler()
	mockUC := h.authUsecase.(*mockAuthUsecase)

	mockUC.On("Logout", mock.Anything).
		Return("expired_token", nil)

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	w := httptest.NewRecorder()

	h.Logout(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp LogoutResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "success", resp.Message)
	assert.Equal(t, "expired_token", resp.Token)
	mockUC.AssertExpectations(t)
}

func TestLogout_Error(t *testing.T) {
	h := newTestHandler()
	mockUC := h.authUsecase.(*mockAuthUsecase)

	mockUC.On("Logout", mock.Anything).
		Return("", errors.New("logout error"))

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	w := httptest.NewRecorder()

	h.Logout(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}
