package handlers

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ----------------------
// validateEmail
// ----------------------
func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"user.name+tag@domain.co", true},
		{"invalid@", false},
		{"@nope.com", false},
		{"", false},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, validateEmail(tt.email), "email: %s", tt.email)
	}
}

// ----------------------
// validatePassword
// ----------------------
func TestValidatePassword(t *testing.T) {
	tests := []struct {
		password string
		expected bool
	}{
		{"Aa1234", true},
		{"abcdef", false}, // нет верхнего регистра, нет цифр
		{"ABCDEF", false}, // нет нижнего регистра, нет цифр
		{"Abcde", false},  // меньше 6 символов
		{"Abcdef", false}, // нет цифр
		{"123456", false}, // нет букв
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, validatePassword(tt.password), "password: %s", tt.password)
	}
}

// ----------------------
// validateLoginRequest
// ----------------------
func TestValidateLoginRequest(t *testing.T) {
	tests := []struct {
		name        string
		req         LoginRequest
		expectError bool
	}{
		{"valid", LoginRequest{"user@example.com", "Password1"}, false},
		{"empty email", LoginRequest{"", "Password1"}, true},
		{"bad email", LoginRequest{"invalid@", "Password1"}, true},
		{"bad password", LoginRequest{"user@example.com", "abc"}, true},
	}

	for _, tt := range tests {
		err := validateLoginRequest(&tt.req)
		if tt.expectError {
			assert.Error(t, err, tt.name)
		} else {
			assert.NoError(t, err, tt.name)
		}
	}
}

// ----------------------
// validateRegisterRequest
// ----------------------
func TestValidateRegisterRequest(t *testing.T) {
	tests := []struct {
		name        string
		req         RegisterRequest
		expectError bool
	}{
		{"valid", RegisterRequest{"user@example.com", "Password1"}, false},
		{"empty email", RegisterRequest{"", "Password1"}, true},
		{"bad email", RegisterRequest{"invalid@", "Password1"}, true},
		{"weak password", RegisterRequest{"user@example.com", "abc"}, true},
	}

	for _, tt := range tests {
		err := validateRegisterRequest(&tt.req)
		if tt.expectError {
			assert.Error(t, err, tt.name)
		} else {
			assert.NoError(t, err, tt.name)
		}
	}
}

// ----------------------
// parseIntQueryParam
// ----------------------
func TestParseIntQueryParam(t *testing.T) {
	req := &http.Request{URL: &url.URL{RawQuery: "page=5"}}

	val, err := parseIntQueryParam(req, "page", 1)
	assert.NoError(t, err)
	assert.Equal(t, 5, val)

	val, err = parseIntQueryParam(req, "missing", 10)
	assert.NoError(t, err)
	assert.Equal(t, 10, val)

	req.URL.RawQuery = "page=abc"
	val, err = parseIntQueryParam(req, "page", 1)
	assert.Error(t, err)
}

// ----------------------
// GetPathParameter
// ----------------------
func TestGetPathParameter(t *testing.T) {
	tests := []struct {
		name        string
		urlPath     string
		basePattern string
		expected    string
	}{
		{"simple", "/api/v1/complexes/123", "/api/v1/complexes/", "123"},
		{"no slash in base", "/api/v1/complexes/123", "/api/v1/complexes", "123"},
		{"extra segments", "/api/v1/complexes/123/edit", "/api/v1/complexes/", "123"},
		{"no match", "/api/v1/other/123", "/api/v1/complexes/", ""},
		{"empty remainder", "/api/v1/complexes/", "/api/v1/complexes/", ""},
	}

	for _, tt := range tests {
		req := &http.Request{URL: &url.URL{Path: tt.urlPath}, Method: http.MethodGet}
		got := GetPathParameter(req, tt.basePattern)
		assert.Equal(t, tt.expected, got, tt.name)
	}
}
