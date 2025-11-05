package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"go.uber.org/zap"
)

// --- Вспомогательная функция для создания multipart-запроса ---
func createMultipartRequest(t *testing.T, fieldName, fileName string, fileContent []byte) (*http.Request, *bytes.Buffer) {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		t.Fatalf("ошибка создания формы: %v", err)
	}

	if _, err := io.Copy(part, bytes.NewReader(fileContent)); err != nil {
		t.Fatalf("ошибка записи контента: %v", err)
	}
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/image/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, body
}

func TestUploadImage_Success(t *testing.T) {
	tmpDir := t.TempDir()
	logger := log.New(zap.NewNop())
	handler := NewImageHandler(logger, "http://localhost:8080", tmpDir)

	// Создаём multipart-запрос
	fileContent := []byte("fake image content")
	req, _ := createMultipartRequest(t, "image", "test.jpg", fileContent)
	rec := httptest.NewRecorder()

	handler.UploadImage(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("ожидался статус 201, получен %d", res.StatusCode)
	}

	var resp map[string]string
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("ошибка декодирования ответа: %v", err)
	}

	url, ok := resp["url"]
	if !ok || url == "" {
		t.Fatalf("в ответе отсутствует поле url: %v", resp)
	}

	// Проверяем, что файл реально создан
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("ошибка чтения директории: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("ожидался 1 сохранённый файл, найдено %d", len(files))
	}
}

func TestUploadImage_BadRequest_NoFile(t *testing.T) {
	tmpDir := t.TempDir()
	logger := log.New(zap.NewNop())
	handler := NewImageHandler(logger, "http://localhost:8080", tmpDir)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/image/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rec := httptest.NewRecorder()
	handler.UploadImage(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("ожидался статус 400, получен %d", res.StatusCode)
	}
}

func TestUploadImage_FailOnCreateFile(t *testing.T) {
	// передаём несуществующий каталог, чтобы вызвать ошибку os.Create
	logger := log.New(zap.NewNop())
	handler := NewImageHandler(logger, "http://localhost:8080", "/invalid/dir")

	fileContent := []byte("fake image content")
	req, _ := createMultipartRequest(t, "image", "test.png", fileContent)
	rec := httptest.NewRecorder()

	handler.UploadImage(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusInternalServerError {
		t.Fatalf("ожидался статус 500, получен %d", res.StatusCode)
	}
}
