package log

import (
	"bytes"
	"context"
	"strings"
	"testing"

	request_id "github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/middleware/request"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// helper: создаёт тестовый логгер, пишущий в буфер
func newTestLogger(buf *bytes.Buffer) *Logger {
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(buf),
		zap.DebugLevel,
	)
	return New(zap.New(core))
}

func TestLogger_InfoWithRequestID(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	ctx := context.WithValue(context.Background(), request_id.RequestIDKey, "req-123")

	logger.Info(ctx, "test message", zap.String("key", "value"))

	out := buf.String()
	if !strings.Contains(out, `"req-123"`) {
		t.Errorf("лог не содержит request_id: %s", out)
	}
	if !strings.Contains(out, `"test message"`) {
		t.Errorf("лог не содержит сообщение: %s", out)
	}
	if !strings.Contains(out, `"key":"value"`) {
		t.Errorf("лог не содержит поле key=value: %s", out)
	}
}

func TestLogger_InfoWithoutContext(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	logger.Info(nil, "no context msg")

	out := buf.String()
	if !strings.Contains(out, `"no context msg"`) {
		t.Errorf("ожидалось сообщение в логе, но его нет: %s", out)
	}
	if strings.Contains(out, `"request_id"`) {
		t.Errorf("не должно быть request_id при отсутствии контекста: %s", out)
	}
}

func TestLogger_WarnAndError(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	ctx := context.WithValue(context.Background(), request_id.RequestIDKey, "req-777")

	logger.Warn(ctx, "warn message")
	logger.Error(ctx, "error message")

	out := buf.String()

	if !strings.Contains(out, `"warn message"`) {
		t.Errorf("ожидалось сообщение warn message: %s", out)
	}
	if !strings.Contains(out, `"error message"`) {
		t.Errorf("ожидалось сообщение error message: %s", out)
	}
	if strings.Count(out, `"req-777"`) != 2 {
		t.Errorf("request_id должен добавляться в оба лога: %s", out)
	}
}

func TestLogger_WithAddsField(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	newLogger := logger.With(zap.String("module", "auth"))
	newLogger.Info(context.Background(), "test with")

	out := buf.String()
	if !strings.Contains(out, `"module":"auth"`) {
		t.Errorf("ожидалось поле module=auth в логе: %s", out)
	}
}

func TestLogger_AddCallerSkip(t *testing.T) {
	// Проверим, что AddCallerSkip(1) не ломает логгер (просто smoke test)
	var buf bytes.Buffer
	zl := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(&buf),
		zap.DebugLevel,
	))
	logger := New(zl)

	logger.Info(context.Background(), "caller skip test")

	if !strings.Contains(buf.String(), "caller skip test") {
		t.Errorf("лог не содержит сообщение caller skip test: %s", buf.String())
	}
}
