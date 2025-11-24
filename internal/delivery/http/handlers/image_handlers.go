package handlers

import (
	_ "context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/response"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"go.uber.org/zap"

	"github.com/google/uuid"
)

type ImageHandler struct {
	logger     *log.Logger
	baseURL    string
	storageDir string
}

const MAX_SIZE = 10 << 20

func NewImageHandler(logger *log.Logger, baseURL, storageDir string) *ImageHandler {
	return &ImageHandler{
		logger:     logger,
		baseURL:    baseURL,
		storageDir: storageDir,
	}
}

// UploadImage — POST /api/v1/image/upload
func (h *ImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(MAX_SIZE); err != nil {
		h.logger.Error(ctx, "failed to parse multipart form", zap.Error(err))
		response.HandleError(w, err, http.StatusBadRequest, "не получилось загрузить картинку")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		h.logger.Error(ctx, "failed to get form file", zap.Error(err))
		response.HandleError(w, err, http.StatusBadRequest, "не получилось загрузить картинку")
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		ext = ".jpg"
	}

	uuidName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	savePath := filepath.Join(h.storageDir, uuidName)

	out, err := os.Create(savePath)
	if err != nil {
		h.logger.Error(ctx, "failed to create file", zap.Error(err))
		response.HandleError(w, err, http.StatusInternalServerError, "не удалось сохранить картинку")
		return
	}
	defer out.Close()

	if _, err = io.Copy(out, file); err != nil {
		h.logger.Error(ctx, "failed to copy file", zap.Error(err))
		response.HandleError(w, err, http.StatusInternalServerError, "не удалось сохранить картинку")
		return
	}

	fileURL := fmt.Sprintf("%s/api/v1/image/%s", h.baseURL, uuidName)

	h.logger.Info(ctx, "image uploaded successfully", zap.String("file", uuidName))
	response.WriteJSON(w, http.StatusCreated, map[string]string{"url": fileURL})
}

func RestrictedImageServer(directory string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        filePath := strings.TrimPrefix(r.URL.Path, "/api/v1/image/")
        if filePath == "" || strings.HasSuffix(filePath, "/") || strings.Contains(filePath, "..") {
            http.Error(w, "Access denied", http.StatusForbidden)
            return
        }

        fullPath := filepath.Join(directory, filePath)
        fileInfo, err := os.Stat(fullPath)
        if os.IsNotExist(err) || fileInfo.IsDir() {
            http.Error(w, "File not found", http.StatusNotFound)
            return
        }

        http.ServeFile(w, r, fullPath)
    })
}