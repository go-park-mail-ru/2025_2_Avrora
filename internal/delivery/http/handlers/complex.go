package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/response"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IComplexUsecase interface {
	GetByID(ctx context.Context, id string) (*domain.HousingComplex, error)
	List(ctx context.Context, page, limit int) (*domain.ComplexesInFeed, error)
	Create(ctx context.Context, complex *domain.HousingComplex) error
	Update(ctx context.Context, complex *domain.HousingComplex) error
	Delete(ctx context.Context, id string) error
}

type ComplexHandler struct {
	complexUsecase IComplexUsecase
	logger         *log.Logger
}

func NewComplexHandler(complexUC IComplexUsecase, logger *log.Logger) *ComplexHandler {
	return &ComplexHandler{
		complexUsecase: complexUC,
		logger:         logger,
	}
}


func (h *ComplexHandler) GetComplexByID(w http.ResponseWriter, r *http.Request) {
	id := GetPathParameter(r, "/api/v1/complexes/")
	if id == "" {
		h.logger.Error(r.Context(), "missing complex ID")
		response.HandleError(w, nil, http.StatusBadRequest, "отсутствует ID жилого комплекса")
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		h.logger.Warn(r.Context(), "invalid UUID format", zap.String("id", id))
		response.HandleError(w, nil, http.StatusBadRequest, "некорректный формат UUID")
		return
	}

	complex, err := h.complexUsecase.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrComplexNotFound) {
			response.HandleError(w, nil, http.StatusNotFound, "жилой комплекс не найден")
			return
		}
		h.logger.Error(r.Context(), "failed to get complex", zap.String("id", id), zap.Error(err))
		response.HandleError(w, nil, http.StatusInternalServerError, "ошибка получения жилого комплекса")
		return
	}

	response.WriteJSON(w, http.StatusOK, complex)
}

// ListComplexes handles GET /api/v1/complexes
func (h *ComplexHandler) ListComplexes(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			response.HandleError(w, nil, http.StatusBadRequest, "некорректный номер страницы")
			return
		}
		page = p
	}

	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l < 1 || l > 100 {
			response.HandleError(w, nil, http.StatusBadRequest, "limit должен быть от 1 до 100")
			return
		}
		limit = l
	}

	result, err := h.complexUsecase.List(r.Context(), page, limit)
	if err != nil {
		h.logger.Error(r.Context(), "failed to list complexes", zap.Error(err))
		response.HandleError(w, nil, http.StatusInternalServerError, "ошибка получения списка жилых комплексов")
		return
	}

	response.WriteJSON(w, http.StatusOK, result)
}

func (h *ComplexHandler) CreateComplex(w http.ResponseWriter, r *http.Request) {
	var req CreateComplexRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(r.Context(), "invalid JSON", zap.Error(err))
		response.HandleError(w, nil, http.StatusBadRequest, "некорректный JSON")
		return
	}

	complexID := uuid.NewString()
	complex := &domain.HousingComplex{
		ID:            complexID,
		Name:          req.Name,
		Description:   req.Description,
		YearBuilt:     req.YearBuilt,
		LocationID:    req.LocationID,
		Developer:     req.Developer,
		Address:       req.Address,
		StartingPrice: req.StartingPrice,
		ImageURLs:     req.ImageURLs,
	}

	err := h.complexUsecase.Create(r.Context(), complex)
	if err != nil {
		response.HandleError(w, nil, http.StatusInternalServerError, "ошибка создания жилого комплекса")
		return
	}

	response.WriteJSON(w, http.StatusCreated, complex)
}

func (h *ComplexHandler) UpdateComplex(w http.ResponseWriter, r *http.Request) {
	id := GetPathParameter(r, "/api/v1/complexes/update/")
	if id == "" {
		response.HandleError(w, nil, http.StatusBadRequest, "отсутствует ID жилого комплекса")
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		response.HandleError(w, nil, http.StatusBadRequest, "некорректный формат UUID")
		return
	}

	var req UpdateComplexRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.HandleError(w, nil, http.StatusBadRequest, "некорректный JSON")
		return
	}

	existing, err := h.complexUsecase.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrComplexNotFound) {
			response.HandleError(w, nil, http.StatusNotFound, "жилой комплекс не найден")
			return
		}
		h.logger.Error(r.Context(), "failed to fetch complex for update", zap.String("id", id), zap.Error(err))
		response.HandleError(w, nil, http.StatusInternalServerError, "ошибка обновления жилого комплекса")
		return
	}

	existing.Name = req.Name
	existing.Description = req.Description
	existing.YearBuilt = req.YearBuilt
	existing.LocationID = req.LocationID
	existing.Developer = req.Developer
	existing.Address = req.Address
	existing.StartingPrice = req.StartingPrice
	existing.ImageURLs = req.ImageURLs

	err = h.complexUsecase.Update(r.Context(), existing)
	if err != nil {
		h.logger.Error(r.Context(), "failed to update complex", zap.String("id", id), zap.Error(err))
		response.HandleError(w, nil, http.StatusInternalServerError, "ошибка обновления жилого комплекса")
		return
	}

	response.WriteJSON(w, http.StatusOK, existing)
}

func (h *ComplexHandler) DeleteComplex(w http.ResponseWriter, r *http.Request) {
	id := GetPathParameter(r, "/api/v1/complexes/delete/")
	if id == "" {
		response.HandleError(w, nil, http.StatusBadRequest, "отсутствует ID жилого комплекса")
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		response.HandleError(w, nil, http.StatusBadRequest, "некорректный формат UUID")
		return
	}

	err := h.complexUsecase.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrComplexNotFound) {
			response.HandleError(w, nil, http.StatusNotFound, "жилой комплекс не найден")
			return
		}
		h.logger.Error(r.Context(), "failed to delete complex", zap.String("id", id), zap.Error(err))
		response.HandleError(w, nil, http.StatusInternalServerError, "ошибка удаления жилого комплекса")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}