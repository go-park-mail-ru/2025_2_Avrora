package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-park-mail-ru/2025_2_Avrora/db"
	"github.com/go-park-mail-ru/2025_2_Avrora/utils"
	_ "github.com/lib/pq"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	dsn := "postgres://postgres:postgres@localhost/2025_2_Avrora_test?sslmode=disable"
	db.InitDB(dsn)
	testDB = db.DB

	// Создаём таблицу, если не существует
	_, err := testDB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password TEXT NOT NULL
		);
	`)
	if err != nil {
		panic("Failed to create table: " + err.Error())
	}

	// Запускаем тесты
	code := m.Run()

	// После тестов — очищаем (опционально)
	// _, _ = testDB.Exec("DROP TABLE IF EXISTS users")

	os.Exit(code)
}

func clearUsersTable() {
	_, _ = testDB.Exec("DELETE FROM users")
}

func TestRegisterHandler_Success(t *testing.T) {
	clearUsersTable()

	body := `{"email": "test@example.com", "password": "secret123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	RegisterHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp AuthResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	if resp.User["email"] != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", resp.User["email"])
	}
	if resp.Token == "" {
		t.Error("Expected JWT token, got empty")
	}
}

func TestRegisterHandler_DuplicateEmail(t *testing.T) {
	clearUsersTable()

	body := `{"email": "duplicate@example.com", "password": "secret123"}`
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBufferString(body))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	RegisterHandler(w1, req1)

	req2 := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	RegisterHandler(w2, req2)

	if w2.Code != http.StatusConflict {
		t.Errorf("Expected status %d, got %d", http.StatusConflict, w2.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w2.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}
	if resp["error"] != "User already exists" {
		t.Errorf("Expected error 'User already exists', got '%s'", resp["error"])
	}
}

func TestLoginHandler_Success(t *testing.T) {
	clearUsersTable()

	// Сначала регистрируем пользователя
	hashedPassword, _ := utils.HashPassword("correct_password")
	_, err := testDB.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", "login@example.com", hashedPassword)
	if err != nil {
		t.Fatal("Failed to insert test user:", err)
	}

	// Пробуем залогиниться
	body := `{"email": "login@example.com", "password": "correct_password"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	LoginHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp AuthResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	if resp.User["email"] != "login@example.com" {
		t.Errorf("Expected email login@example.com, got %s", resp.User["email"])
	}
	if resp.Token == "" {
		t.Error("Expected JWT token, got empty")
	}
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	clearUsersTable()

	// Регистрируем пользователя
	hashedPassword, _ := utils.HashPassword("correct_password")
	_, err := testDB.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", "login@example.com", hashedPassword)
	if err != nil {
		t.Fatal("Failed to insert test user:", err)
	}

	// Неверный пароль
	body := `{"email": "login@example.com", "password": "wrong_password"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	LoginHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}
	if resp["error"] != "Invalid credentials" {
		t.Errorf("Expected error 'Invalid credentials', got '%s'", resp["error"])
	}
}

func TestLoginHandler_UserNotFound(t *testing.T) {
	clearUsersTable()

	body := `{"email": "notfound@example.com", "password": "any_password"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	LoginHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}
	if resp["error"] != "Invalid credentials" {
		t.Errorf("Expected error 'Invalid credentials', got '%s'", resp["error"])
	}
}

func TestLogoutHandler(t *testing.T) {
	clearUsersTable()

	// Регистрируем пользователя
	hashedPassword, _ := utils.HashPassword("correct_password")
	_, err := testDB.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", "login@example.com", hashedPassword)
	if err != nil {
		t.Fatal("Failed to insert test user:", err)
	}

	// Пытаемся выйти
	body := `{"email": "login@example.com", "password": "correct_password"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	LogoutHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp map[string]string

	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}
	if resp["message"] != "Logged out successfully" {
		t.Errorf("Expected message 'Logged out successfully', got '%s'", resp["message"])
	}
	if resp["Token"] == "" {
		t.Error("Expected JWT token, got empty")
	}
}

func TestGenerateInvalidToken(t *testing.T) {
	_, err := generateExpiredJWT()
	if err != nil {
		t.Errorf("Expected error, got nil")
	}
}