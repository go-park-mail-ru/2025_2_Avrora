package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- Мок интерфейса profileUsecase ---
type mockProfileUsecase struct {
	getProfileByIDFunc            func(ctx context.Context, id string) (*domain.Profile, error)
	updateProfileFunc             func(ctx context.Context, id string, update *domain.ProfileUpdate) error
	updateProfileSecurityByIDFunc func(ctx context.Context, id, oldPass, newPass string) error
	updateEmailFunc               func(ctx context.Context, id, email string) error
}

func (m *mockProfileUsecase) GetProfileByID(ctx context.Context, id string) (*domain.Profile, error) {
	return m.getProfileByIDFunc(ctx, id)
}

func (m *mockProfileUsecase) UpdateProfile(ctx context.Context, id string, update *domain.ProfileUpdate) error {
	return m.updateProfileFunc(ctx, id, update)
}

func (m *mockProfileUsecase) UpdateProfileSecurityByID(ctx context.Context, id, oldPass, newPass string) error {
	return m.updateProfileSecurityByIDFunc(ctx, id, oldPass, newPass)
}

func (m *mockProfileUsecase) UpdateEmail(ctx context.Context, id, email string) error {
	return m.updateEmailFunc(ctx, id, email)
}

// --- Вспомогательная функция для создания логгера ---
func newTestLogger() *log.Logger {
	core, _ := observer.New(zapcore.InfoLevel)
	zapLogger := zap.New(core)
	return log.New(zapLogger)
}

// --- Тесты ---
func TestProfileHandler_GetProfile(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		mockReturn *domain.Profile
		mockErr    error
		wantStatus int
	}{
		{
			name:       "успешное получение профиля",
			url:        "/api/v1/profile/123",
			mockReturn: &domain.Profile{ID: "123", FirstName: "Test"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "ошибка получения профиля",
			url:        "/api/v1/profile/123",
			mockErr:    errors.New("db error"),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "пустой id",
			url:        "/api/v1/profile/",
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := &mockProfileUsecase{
				getProfileByIDFunc: func(ctx context.Context, id string) (*domain.Profile, error) {
					return tt.mockReturn, tt.mockErr
				},
			}

			handler := &profileHandler{
				profileUsecase: mockUC,
				log:            newTestLogger(),
			}

			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()

			handler.GetProfile(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("ожидался статус %d, получен %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestProfileHandler_UpdateProfile(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		body       interface{}
		mockErr    error
		wantStatus int
	}{
		{
			name:       "успешное обновление",
			url:        "/api/v1/profile/update/123",
			body:       ProfileUpdate{FirstName: ptr("Ivan")},
			wantStatus: http.StatusOK,
		},
		{
			name:       "невалидный JSON",
			url:        "/api/v1/profile/update/123",
			body:       "{invalid-json}",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "пустой id",
			url:        "/api/v1/profile/update/",
			body:       ProfileUpdate{FirstName: ptr("Ivan")},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "ошибка usecase",
			url:        "/api/v1/profile/update/123",
			body:       ProfileUpdate{FirstName: ptr("Ivan")},
			mockErr:    errors.New("db error"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := &mockProfileUsecase{
				updateProfileFunc: func(ctx context.Context, id string, update *domain.ProfileUpdate) error {
					return tt.mockErr
				},
			}

			handler := &profileHandler{
				profileUsecase: mockUC,
				log:            newTestLogger(),
			}

			var bodyBytes []byte
			switch v := tt.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				bodyBytes, _ = json.Marshal(v)
			}

			req := httptest.NewRequest(http.MethodPost, tt.url, bytes.NewBuffer(bodyBytes))
			w := httptest.NewRecorder()

			handler.UpdateProfile(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("ожидался статус %d, получен %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestProfileHandler_UpdateProfileSecurityByID(t *testing.T) {
	mockUC := &mockProfileUsecase{
		updateProfileSecurityByIDFunc: func(ctx context.Context, id, oldPass, newPass string) error {
			if oldPass == "wrong" {
				return errors.New("bad password")
			}
			return nil
		},
	}
	handler := &profileHandler{profileUsecase: mockUC, log: newTestLogger()}

	// --- успешное обновление
	reqBody := `{"OldPassword":"old","NewPassword":"new"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/profile/security/123", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()
	handler.UpdateProfileSecurityByID(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("ожидался 200, получен %d", w.Code)
	}

	// --- невалидный JSON
	req = httptest.NewRequest(http.MethodPost, "/api/v1/profile/security/123", bytes.NewBufferString("{bad-json}"))
	w = httptest.NewRecorder()
	handler.UpdateProfileSecurityByID(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("ожидался 400, получен %d", w.Code)
	}
}

func TestProfileHandler_UpdateEmail(t *testing.T) {
	mockUC := &mockProfileUsecase{
		updateEmailFunc: func(ctx context.Context, id, email string) error {
			if email == "bad@example.com" {
				return errors.New("db error")
			}
			return nil
		},
	}
	handler := &profileHandler{profileUsecase: mockUC, log: newTestLogger()}

	tests := []struct {
		name       string
		url        string
		body       string
		wantStatus int
	}{
		{"успешное обновление", "/api/v1/profile/email/123", `{"Email":"test@example.com"}`, http.StatusOK},
		{"ошибка usecase", "/api/v1/profile/email/123", `{"Email":"bad@example.com"}`, http.StatusInternalServerError},
		{"невалидный JSON", "/api/v1/profile/email/123", `{bad}`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.url, bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()
			handler.UpdateEmail(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("ожидался %d, получен %d", tt.wantStatus, w.Code)
			}
		})
	}
}

// helper
func ptr(s string) *string { return &s }
