package handlers

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// создаём мок реализации IAuthUsecase (пустую, т.к. не вызывается)
type u struct{}

func (m *u) Register(ctx context.Context, email string, password string) error {
	return nil
}

func (m *u) Login(ctx context.Context, email string, password string) (string, error) {
	return "mock_token", nil
}

func (m *u) Logout(ctx context.Context) (string, error) {
	return "expired_token", nil
}

// ----------------------------------------------
// ✅ Тест конструктора NewAuthHandler
// ----------------------------------------------
func TestNewAuthHandler(t *testing.T) {
	// Подготовка
	mockUC := &u{}
	zapLogger := zap.NewNop()
	appLogger := log.New(zapLogger)

	// Действие
	handler := NewAuthHandler(mockUC, appLogger)

	// Проверки
	assert.NotNil(t, handler, "handler должен быть создан")
	assert.Equal(t, mockUC, handler.authUsecase, "authUsecase должен быть корректно присвоен")
	assert.Equal(t, appLogger, handler.logger, "logger должен быть корректно присвоен")
}
