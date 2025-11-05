package handlers

import (
	"testing"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"go.uber.org/zap"
)

// --- Тест конструктора ---
func TestNewProfileHandler(t *testing.T) {
	mockUC := &mockProfileUsecase{}
	logger := log.New(zap.NewNop())

	handler := NewProfileHandler(mockUC, logger)

	if handler == nil {
		t.Fatal("ожидался непустой handler, получен nil")
	}
	if handler.profileUsecase != mockUC {
		t.Error("profileUsecase установлен неверно")
	}
	if handler.log != logger {
		t.Error("логгер установлен неверно")
	}
}

// --- Тест интерфейсной совместимости ---
func TestProfileHandler_ImplementsInterface(t *testing.T) {
	var _ IProfileUsecase = &mockProfileUsecase{} // компилятор проверит совместимость
}
