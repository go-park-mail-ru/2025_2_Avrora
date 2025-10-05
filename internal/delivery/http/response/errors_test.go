package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewErrorResp(t *testing.T) {
	msg := "ошибка"
	resp := NewErrorResp(msg)
	if resp.Error != msg {
		t.Errorf("expected %q, got %q", msg, resp.Error)
	}
}

func TestWriteJSON_Basic(t *testing.T) {
	rec := httptest.NewRecorder()
	data := map[string]string{"key": "value"}

	WriteJSON(rec, http.StatusOK, data)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", rec.Header().Get("Content-Type"))
	}

	var decoded map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&decoded); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if decoded["key"] != "value" {
		t.Errorf("expected value 'value', got '%s'", decoded["key"])
	}
}

func TestHandleError_Basic(t *testing.T) {
	rec := httptest.NewRecorder()
	HandleError(rec, nil, http.StatusBadRequest, "ошибка пользователя")

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var resp ErrorResp
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error != "ошибка пользователя" {
		t.Errorf("expected error 'ошибка пользователя', got '%s'", resp.Error)
	}
}
