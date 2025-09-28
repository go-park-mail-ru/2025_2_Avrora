package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-park-mail-ru/2025_2_Avrora/db"
	"github.com/go-park-mail-ru/2025_2_Avrora/models"
	"github.com/go-park-mail-ru/2025_2_Avrora/response"
	"github.com/go-park-mail-ru/2025_2_Avrora/utils"
)

func setupTestRepo(t *testing.T) *db.Repo {
	t.Helper()

	// задаём окружение прямо здесь → тесты независимы от .env
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "testdb") // всегда отдельная тестовая база
	os.Setenv("PASSWORD_PEPPER", "pepper123")

	dsn := utils.GetPostgresDSN()
	repo, err := db.New(dsn)
	if err != nil {
		t.Fatal("failed to connect to test DB:", err)
	}
	if err := repo.ClearAllTables(); err != nil {
		t.Fatal("failed to clear tables:", err)
	}
	return repo
}

func TestRegisterHandler_Success(t *testing.T) {
	jwtGen := utils.NewJwtGenerator("test_secret_32_chars_min_for_tests")
	repo := setupTestRepo(t)
	passwordHasher, _ := utils.NewPasswordHasher(os.Getenv("PASSWORD_PEPPER"))

	body := `{"email": "test@example.com", "password": "secret123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	RegisterHandler(w, req, repo, jwtGen, passwordHasher)

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
	jwtGen := utils.NewJwtGenerator("test_secret_32_chars_min_for_tests")
	repo := setupTestRepo(t)
	passwordHasher, _ := utils.NewPasswordHasher(os.Getenv("PASSWORD_PEPPER"))

	body := `{"email": "duplicate@example.com", "password": "secret123"}`
	// первый запрос
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBufferString(body))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	RegisterHandler(w1, req1, repo, jwtGen, passwordHasher)

	// второй запрос (тот же email)
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBufferString(body))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	RegisterHandler(w2, req2, repo, jwtGen, passwordHasher)

	if w2.Code != http.StatusConflict {
		t.Errorf("Expected status %d, got %d", http.StatusConflict, w2.Code)
	}

	var resp response.ErrorResp
	if err := json.NewDecoder(w2.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}
	if resp.Error != "Пользователь с таким email уже существует" {
		t.Errorf("Expected error message, got '%s'", resp.Error)
	}
}

func TestLoginHandler_Success(t *testing.T) {
	jwtGen := utils.NewJwtGenerator("test_secret_32_chars_min_for_tests")
	repo := setupTestRepo(t)
	passwordHasher, _ := utils.NewPasswordHasher(os.Getenv("PASSWORD_PEPPER"))

	// регаем пользователя вручную
	hashedPassword, _ := passwordHasher.Hash("correct_password")
	user := models.User{Email: "login@example.com", Password: hashedPassword}
	if err := repo.User().Create(&user); err != nil {
		t.Fatal("Failed to insert test user:", err)
	}

	body := `{"email": "login@example.com", "password": "correct_password"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	LoginHandler(w, req, repo, jwtGen, passwordHasher)

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
	jwtGen := utils.NewJwtGenerator("test_secret_32_chars_min_for_tests")
	repo := setupTestRepo(t)
	passwordHasher, _ := utils.NewPasswordHasher(os.Getenv("PASSWORD_PEPPER"))

	// регаем пользователя
	hashedPassword, _ := passwordHasher.Hash("correct_password")
	user := models.User{Email: "login@example.com", Password: hashedPassword}
	if err := repo.User().Create(&user); err != nil {
		t.Fatal("Failed to insert test user:", err)
	}

	// пробуем неверный пароль
	body := `{"email": "login@example.com", "password": "wrong_password"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	LoginHandler(w, req, repo, jwtGen, passwordHasher)

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
	jwtGen := utils.NewJwtGenerator("test_secret_32_chars_min_for_tests")
	repo := setupTestRepo(t)
	passwordHasher, _ := utils.NewPasswordHasher(os.Getenv("PASSWORD_PEPPER"))

	hashedPassword, _ := passwordHasher.Hash("correct_password")
	user := models.User{Email: "login@example.com", Password: hashedPassword}
	if err := repo.User().Create(&user); err != nil {
		t.Fatal("Failed to insert test user:", err)
	}

	// логинимся
	loginBody := `{"email": "login@example.com", "password": "correct_password"}`
	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	LoginHandler(loginW, loginReq, repo, jwtGen, passwordHasher)

	var loginResp AuthResponse
	if err := json.NewDecoder(loginW.Body).Decode(&loginResp); err != nil {
		t.Fatal("Failed to decode login response:", err)
	}

	// logout
	logoutReq := httptest.NewRequest(http.MethodPost, "/api/v1/logout", nil)
	logoutReq.Header.Set("Authorization", "Bearer "+loginResp.Token)
	logoutW := httptest.NewRecorder()

	LogoutHandler(logoutW, logoutReq, jwtGen)

	if logoutW.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, logoutW.Code)
	}

	var resp struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(logoutW.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}
	if resp.Message != "успешный логаут" {
		t.Errorf("Expected message 'успешный логаут', got '%s'", resp.Message)
	}
}
