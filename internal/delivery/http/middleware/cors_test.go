package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot) // 418 для проверки что дошло
	w.Write([]byte("ok"))
}

func TestCorsMiddleware_SetsHeaders(t *testing.T) {
	handler := CorsMiddleware(http.HandlerFunc(dummyHandler), "http://example.com")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	tests := map[string]string{
		"Access-Control-Allow-Origin":      "http://example.com",
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, OPTIONS",
		"Access-Control-Allow-Headers":     "Content-Type, Authorization",
	}

	for key, expected := range tests {
		if got := res.Header.Get(key); got != expected {
			t.Errorf("ожидался заголовок %s=%s, получен %s", key, expected, got)
		}
	}

	if res.StatusCode != http.StatusTeapot {
		t.Errorf("ожидался код 418, получен %d", res.StatusCode)
	}
}

func TestCorsMiddleware_OptionsRequest(t *testing.T) {
	handler := CorsMiddleware(http.HandlerFunc(dummyHandler), "*")

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("ожидался статус 200 для OPTIONS, получен %d", res.StatusCode)
	}

	// Проверим, что dummyHandler не вызывался (тело пустое)
	if rec.Body.Len() != 0 {
		t.Errorf("ожидалось пустое тело для OPTIONS, получено: %q", rec.Body.String())
	}
}
