package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ImageHandler struct {
	logger     *zap.Logger
	baseURL    string
	storageDir string
}

func NewImageHandler(logger *zap.Logger, baseURL, storageDir string) *ImageHandler {
	return &ImageHandler{
		logger:     logger,
		baseURL:    baseURL,
		storageDir: storageDir,
	}
}

// UploadImage — POST /api/v1/image/upload
func (h *ImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // до 10 MB
	if err != nil {
		h.logger.Error("failed to parse multipart form", zap.Error(err))
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		h.logger.Error("failed to get form file", zap.Error(err))
		http.Error(w, "no file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Получаем расширение исходного файла
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		ext = ".jpg" // можно задать дефолтное расширение
	}

	// Генерируем UUID и формируем новое имя файла
	uuidName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	savePath := filepath.Join(h.storageDir, uuidName)

	out, err := os.Create(savePath)
	if err != nil {
		h.logger.Error("failed to create file", zap.Error(err))
		http.Error(w, "could not save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	if _, err = io.Copy(out, file); err != nil {
		h.logger.Error("failed to copy file", zap.Error(err))
		http.Error(w, "could not save file", http.StatusInternalServerError)
		return
	}

	fileURL := fmt.Sprintf("%s/api/v1/image/%s", h.baseURL, uuidName)

	h.logger.Info("image uploaded successfully", zap.String("file", uuidName))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"url": "%s"}`, fileURL)
}
