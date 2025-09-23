package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"

	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-park-mail-ru/2025_2_Avrora/db"
	"github.com/go-park-mail-ru/2025_2_Avrora/models"
	"github.com/go-park-mail-ru/2025_2_Avrora/response"
	"github.com/go-park-mail-ru/2025_2_Avrora/utils"
)

var testRepo *db.Repo

func TestMain(m *testing.M) {
	utils.LoadEnv()
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)
	testRepo = db.NewRepo()
	if err := testRepo.Init(dsn); err != nil {
		panic("Failed to initialize test DB: " + err.Error())
	}

	testRepo.User().ClearUserTable()

	code := m.Run()

	os.Exit(code)
}

func TestRegisterHandler_Success(t *testing.T) {
	body := `{"email": "test@example.com", "password": "secret123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	RegisterHandler(w, req, testRepo)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp AuthResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	if resp.User.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", resp.User.Email)
	}
	if resp.Token == "" {
		t.Error("Expected JWT token, got empty")
	}
}

func TestRegisterHandler_DuplicateEmail(t *testing.T) {
	// Регистрируем первый раз
	body := `{"email": "duplicate@example.com", "password": "secret123!В"}`
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBufferString(body))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	RegisterHandler(w1, req1, testRepo)

	// Повторная регистрация
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBufferString(body))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	RegisterHandler(w2, req2, testRepo)

	if w2.Code != http.StatusConflict {
		t.Errorf("Expected status %d, got %d", http.StatusConflict, w2.Code)
	}

	var resp response.ErrorResp
	if err := json.NewDecoder(w2.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}
	if resp.Error != "Пользователь с таким email уже существует" {
		t.Errorf("Expected error 'Пользователь с таким email уже существует', got '%s'", resp.Error)
	}
}

func TestLoginHandler_Success(t *testing.T) {
	// Сначала регистрируем пользователя
	hashedPassword, _ := utils.HashPassword("correct_pasВ3sword!")
	user := models.User{Email: "login@example.com", Password: hashedPassword}
	if err := testRepo.User().Create(&user); err != nil {
		t.Fatal("Failed to insert test user:", err)
	}

	// Пробуем залогиниться
	body := `{"email": "login@example.com", "password": "correct_pasВ3sword!"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	LoginHandler(w, req, testRepo)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp AuthResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	if resp.User.Email != "login@example.com" {
		t.Errorf("Expected email login@example.com, got %s", resp.User.Email)
	}
	if resp.Token == "" {
		t.Error("Expected JWT token, got empty")
	}
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	testRepo.User().ClearUserTable()
	// Регистрируем пользователя
	hashedPassword, _ := utils.HashPassword("correct_password!В3")
	user := models.User{Email: "login@example.com", Password: hashedPassword}
	if err := testRepo.User().Create(&user); err != nil {
		t.Fatal("Failed to insert test user:", err)
	}

	// Неверный пароль
	body := `{"email": "login@example.com", "password": "wrong_password!В3"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	LoginHandler(w, req, testRepo)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var resp response.ErrorResp
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}
	if resp.Error != "невалидные учетные данные" {
		t.Errorf("Expected error 'невалидные учетные данные', got '%s'", resp.Error)
	}
}

func TestLogoutHandler(t *testing.T) {
	testRepo.User().ClearUserTable()
	hashedPassword, _ := utils.HashPassword("correct_password")
	user := models.User{Email: "login@example.com", Password: hashedPassword}
	if err := testRepo.User().Create(&user); err != nil {
		t.Fatal("Failed to insert test user:", err)
	}

	// Логинимся, чтобы получить токен
	loginBody := `{"email": "login@example.com", "password": "correct_password"}`
	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	LoginHandler(loginW, loginReq, testRepo)

	var loginResp AuthResponse
	if err := json.NewDecoder(loginW.Body).Decode(&loginResp); err != nil {
		t.Fatal("Failed to decode login response:", err)
	}

	// Выполняем logout
	logoutReq := httptest.NewRequest(http.MethodPost, "/api/v1/logout", nil)
	logoutReq.Header.Set("Authorization", "Bearer "+loginResp.Token)
	logoutW := httptest.NewRecorder()

	LogoutHandler(logoutW, logoutReq)

	if logoutW.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, logoutW.Code)
	}

	var resp struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(logoutW.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	if resp.Message != "успещный логаут" {
		t.Errorf("Expected message 'успещный логаут', got '%s'", resp.Message)
	}
}