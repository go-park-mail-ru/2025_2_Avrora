package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCorsMiddleware(t *testing.T) {
	corsOrigin := "http://example.com"

	// Заглушка следующего handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot) // чтобы проверить, что вызвали next
	})

	// Создаем middleware
	handler := CorsMiddleware(nextHandler, corsOrigin)

	// === Тест OPTIONS (preflight) ===
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for OPTIONS, got %d", resp.StatusCode)
	}
	if resp.Header.Get("Access-Control-Allow-Origin") != corsOrigin {
		t.Errorf("Expected Access-Control-Allow-Origin %s, got %s", corsOrigin, resp.Header.Get("Access-Control-Allow-Origin"))
	}

	// === Тест GET (должен вызвать next) ===
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	w = httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp = w.Result()
	if resp.StatusCode != http.StatusTeapot {
		t.Errorf("Expected status 418 from next handler, got %d", resp.StatusCode)
	}
	if resp.Header.Get("Access-Control-Allow-Origin") != corsOrigin {
		t.Errorf("Expected Access-Control-Allow-Origin %s, got %s", corsOrigin, resp.Header.Get("Access-Control-Allow-Origin"))
	}
}
