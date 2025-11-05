package middleware

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	request_id "github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/middleware/request"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// --- вспомогательный логгер, пишущий в bytes.Buffer ---
func newTestLogger(buf *bytes.Buffer) *log.Logger {
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(buf),
		zap.DebugLevel,
	)
	return log.New(zap.New(core))
}

// --- простой handler, возвращающий статус ---
func handlerWithStatus(code int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond) // имитация работы
		w.WriteHeader(code)
	}
}

// --- тесты ---
func TestLoggerMiddleware_InfoLog(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	middleware := LoggerMiddleware(logger)
	handler := middleware(handlerWithStatus(http.StatusOK))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req = req.WithContext(context.WithValue(req.Context(), request_id.RequestIDKey, "req-123"))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("ожидался 200, получен %d", rec.Code)
	}

	out := buf.String()
	if !strings.Contains(out, `"request completed"`) {
		t.Errorf("лог не содержит 'request completed': %s", out)
	}
	if !strings.Contains(out, `"request_id":"req-123"`) {
		t.Errorf("лог не содержит request_id: %s", out)
	}
}

func TestLoggerMiddleware_WarnLog(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	middleware := LoggerMiddleware(logger)
	handler := middleware(handlerWithStatus(http.StatusBadRequest))

	req := httptest.NewRequest(http.MethodGet, "/warn", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("ожидался 400, получен %d", rec.Code)
	}

	out := buf.String()
	if !strings.Contains(out, `"client error"`) {
		t.Errorf("лог не содержит 'client error': %s", out)
	}
}

func TestLoggerMiddleware_ErrorLog(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	middleware := LoggerMiddleware(logger)
	handler := middleware(handlerWithStatus(http.StatusInternalServerError))

	req := httptest.NewRequest(http.MethodPost, "/fail", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("ожидался 500, получен %d", rec.Code)
	}

	out := buf.String()
	if !strings.Contains(out, `"server error"`) {
		t.Errorf("лог не содержит 'server error': %s", out)
	}
}

func TestGetClientIP(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("X-Forwarded-For", "192.168.0.1, 10.0.0.1")
	ip := getClientIP(r)
	if ip != "192.168.0.1" {
		t.Errorf("ожидался 192.168.0.1, получен %s", ip)
	}

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Del("X-Forwarded-For")
	r.Header.Set("X-Real-IP", "10.10.10.10")
	ip = getClientIP(r)
	if ip != "10.10.10.10" {
		t.Errorf("ожидался 10.10.10.10, получен %s", ip)
	}

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.RemoteAddr = "127.0.0.1:12345"
	ip = getClientIP(r)
	if ip != "127.0.0.1" {
		t.Errorf("ожидался 127.0.0.1, получен %s", ip)
	}
}
