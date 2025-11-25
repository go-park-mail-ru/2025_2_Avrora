package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/middleware"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/response"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
)

type IOfferUsecase interface {
	ListOffersInFeed(ctx context.Context, page, limit int) (*domain.OffersInFeed, error)
	Get(ctx context.Context, id string) (*domain.Offer, error)
	Update(ctx context.Context, offer *domain.Offer) error
	Create(ctx context.Context, offer *domain.Offer) error
	Delete(ctx context.Context, id string) error
	ListOffersInFeedByUserID(ctx context.Context, userID string, page, limit int) (*domain.OffersInFeed, error)
	FilterOffers(ctx context.Context, f *domain.OfferFilter, limit, offset int) ([]domain.OfferInFeed, error)
	ToggleLike(ctx context.Context, userID, offerID string) (bool, error)
	IsLiked(ctx context.Context, userID, offerID string) (bool, error)
	GetLikesCount(ctx context.Context, offerID string) (int, error)
}

type offerHandler struct {
	offerUsecase IOfferUsecase
	logger       *log.Logger
}

func NewOfferHandler(uc IOfferUsecase, logger *log.Logger) *offerHandler {
	return &offerHandler{offerUsecase: uc, logger: logger}
}
func (h *offerHandler) FilterOffers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	f := &domain.OfferFilter{}

	if v := q.Get("offer_type"); v != "" {
		f.OfferType = &v
	}
	if v := q.Get("property_type"); v != "" {
		f.PropertyType = &v
	}
	if v := q.Get("status"); v != "" {
		f.Status = &v
	}
	if v := q.Get("rooms"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			f.Rooms = &i
		}
	}
	if v := q.Get("price_min"); v != "" {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.PriceMin = &i
		}
	}
	if v := q.Get("price_max"); v != "" {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.PriceMax = &i
		}
	}
	if v := q.Get("area_min"); v != "" {
		if f64, err := strconv.ParseFloat(v, 64); err == nil {
			f.AreaMin = &f64
		}
	}
	if v := q.Get("area_max"); v != "" {
		if f64, err := strconv.ParseFloat(v, 64); err == nil {
			f.AreaMax = &f64
		}
	}
	if v := q.Get("address"); v != "" {
		f.Address = &v
	}

	limit := 20
	offset := 0
	if v := q.Get("limit"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			limit = i
		}
	}
	if v := q.Get("offset"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			offset = i
		}
	}

	offers, err := h.offerUsecase.FilterOffers(ctx, f, limit, offset)
	if err != nil {
		h.logger.Error(ctx, "failed to filter offers", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(offers); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// ToggleLike устаналивает или удаляет лайк у объявления пост запросом
func (h *offerHandler) ToggleLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.HandleError(w, nil, http.StatusMethodNotAllowed, "метод не поддерживается")
		return
	}

	ctx := r.Context()
	offerID := r.URL.Query().Get("id")
	if offerID == "" {
		response.HandleError(w, nil, http.StatusBadRequest, "параметр 'id' обязателен")
		return
	}

	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok || userID == "" {
		response.HandleError(w, nil, http.StatusUnauthorized, "пользователь не авторизован")
		return
	}

	liked, err := h.offerUsecase.ToggleLike(ctx, userID, offerID)
	if err != nil {
		if errors.Is(err, domain.ErrOfferNotFound) {
			response.HandleError(w, err, http.StatusNotFound, "объявление не найдено")
			return
		}
		h.logger.Error(ctx, "failed to toggle like", zap.String("offer_id", offerID), zap.Error(err))
		response.HandleError(w, err, http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"liked":   liked,
		"message": "лайк обновлён",
	})
}

func (h *offerHandler) IsLiked(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.HandleError(w, nil, http.StatusMethodNotAllowed, "метод не поддерживается")
		return
	}

	ctx := r.Context()
	offerID := r.URL.Query().Get("id")
	if offerID == "" {
		response.HandleError(w, nil, http.StatusBadRequest, "параметр 'id' обязателен")
		return
	}

	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok || userID == "" {
		response.WriteJSON(w, http.StatusOK, map[string]bool{"liked": false})
		return
	}

	liked, err := h.offerUsecase.IsLiked(ctx, userID, offerID)
	if err != nil {
		if errors.Is(err, domain.ErrOfferNotFound) {
			response.HandleError(w, err, http.StatusNotFound, "объявление не найдено")
			return
		}
		h.logger.Warn(ctx, "failed to check like status", zap.String("offer_id", offerID), zap.Error(err))
		liked = false
	}

	response.WriteJSON(w, http.StatusOK, map[string]bool{"liked": liked})
}

//тут должен быть эндпоинт для просмотров (делается)
