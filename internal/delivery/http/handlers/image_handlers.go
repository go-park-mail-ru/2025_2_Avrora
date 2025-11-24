package handlers

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/response"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	fileserverpb "github.com/go-park-mail-ru/2025_2_Avrora/proto/fileserver"
	"go.uber.org/zap"

	"github.com/google/uuid"
)

type ImageHandler struct {
	fileserver fileserverpb.FileServerClient // gRPC client for file operations
	logger     *log.Logger
	baseURL    string
}

const MAX_SIZE = 10 << 20 // 10MB

func NewImageHandler(fs fileserverpb.FileServerClient, logger *log.Logger, baseURL string) *ImageHandler {
	return &ImageHandler{
		fileserver: fs,
		logger:     logger,
		baseURL:    baseURL,
	}
}

// UploadImage — POST /api/v1/image/upload
func (h *ImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(MAX_SIZE); err != nil {
		h.logger.Error(ctx, "failed to parse multipart form", zap.Error(err))
		response.HandleError(w, err, http.StatusBadRequest, "размер файла превышает допустимый лимит(10MB)")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		h.logger.Error(ctx, "failed to get form file", zap.Error(err))
		response.HandleError(w, err, http.StatusBadRequest, "не получилось загрузить картинку")
		return
	}
	defer file.Close()

	// Read file data into memory (safe for 10MB limit)
	data, err := io.ReadAll(file)
	if err != nil {
		h.logger.Error(ctx, "failed to read file data", zap.Error(err))
		response.HandleError(w, err, http.StatusInternalServerError, "не удалось прочитать файл")
		return
	}

	// Get and sanitize file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		ext = ".jpg"
	}

	// Generate UUID-based filename
	uuidName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Get content type (fallback if not provided)
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// ✅ Call gRPC file server to upload
	req := &fileserverpb.UploadRequest{
		Data:        data,
		Filename:    uuidName,
		ContentType: contentType,
	}

	resp, err := h.fileserver.Upload(ctx, req)
	if err != nil {
		h.logger.Error(ctx, "gRPC upload failed", zap.Error(err), zap.String("filename", uuidName))
		response.HandleError(w, err, http.StatusInternalServerError, "не удалось сохранить картинку на сервере")
		return
	}

	// Construct full URL (handle both relative and absolute paths from server)
	fileURL := resp.Url
	if !strings.HasPrefix(fileURL, "http") {
		fileURL = fmt.Sprintf("%s%s", h.baseURL, fileURL)
	}

	h.logger.Info(ctx, "image uploaded successfully via gRPC",
		zap.String("original_filename", header.Filename),
		zap.String("stored_filename", uuidName),
		zap.String("url", fileURL),
	)

	response.WriteJSON(w, http.StatusCreated, map[string]string{"url": fileURL})
}

// GetImage — GET /api/v1/image/{filename}
func (h *ImageHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filename := strings.TrimPrefix(r.URL.Path, "/api/v1/image/")
	if filename == "" || strings.Contains(filename, "/") || strings.Contains(filename, "..") {
		response.HandleError(w, nil, http.StatusNotFound, "некорректное имя файла")
		return
	}

	req := &fileserverpb.GetRequest{
		Filename: filename,
	}

	stream, err := h.fileserver.Get(ctx, req)
	if err != nil {
		h.logger.Error(ctx, "gRPC get failed", zap.Error(err), zap.String("filename", filename))
		response.HandleError(w, err, http.StatusNotFound, "файл не найден")
		return
	}

	// Stream response chunks to client
	w.Header().Set("Cache-Control", "public, max-age=86400") // Cache for 24 hours

	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			h.logger.Error(ctx, "stream recv failed", zap.Error(err))
			response.HandleError(w, err, http.StatusInternalServerError, "ошибка при чтении файла")
			return
		}

		if _, writeErr := w.Write(resp.Chunk); writeErr != nil {
			// Client disconnected, we can stop
			return
		}
	}
}

// ImageServer returns a handler that serves images via gRPC
func (h *ImageHandler) ImageServer() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.GetImage(w, r)
	})
}